// Go Skeleton API.
//
// Sample REST API.
//
//     Schemes: http
//     Host: localhost:8081
//     BasePath: /
//     Version: v1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package gin

import (
"fmt"

"github.com/gin-gonic/gin"
"github.com/uniplaces/stripe-gateway/domain"
"github.com/uniplaces/stripe-gateway/gin/handlers"
"github.com/uniplaces/stripe-gateway/gin/handlers/v1/hello"
"github.com/uniplaces/stripe-gateway/gin/middleware"
)

// APIService is a Gin API service
type APIService struct {
	*gin.Engine

	addr string
}

// New creates an initialized Gin APIService
func New(
	host,
	port,
	corsDomain,
	mode string,
	interactorAggregator domain.InteractorAggregate,
) (
	APIService,
	error,
) {
	loggerMiddleware, err := middleware.Logger()
	if err != nil {
		return APIService{}, err
	}

	gin.SetMode(mode)
	router := gin.New()
	router.Use(middleware.Cors(corsDomain), loggerMiddleware)

	sourceHandler := source.NewHandler(interactorAggregator.Source)
	customerHandler := customer.NewHandler(interactorAggregator.Customer)
	chargeHandler := charge.NewHandler(interactorAggregator.Charge)
	accountHandler := account.NewHandler(interactorAggregator.Account)
	eventHandler := event.NewHandler(interactorAggregator.Event)

	router.GET("/ping", handlers.Ping())

	gv1 := router.Group("/v1")
	{
		customers := gv1.Group("/customers")
		{
			customers.GET("/:userId", customerHandler.GetCustomer())
			customers.POST("", customerHandler.Create())
			customers.PUT("/:userId/card", customerHandler.AddCard())
			customers.GET("/:userId/default-card", customerHandler.GetDefaultCard())
			customers.POST("/:userId/charge", customerHandler.Charge())
		}

		accounts := gv1.Group("/accounts")
		{
			accounts.GET("/:userId", accountHandler.GetAccount())
			accounts.POST("/:userId/transfer", accountHandler.Transfer())
			accounts.POST("/:userId/payout", accountHandler.Payout())
		}

		sources := gv1.Group("/sources")
		{
			sources.POST("", sourceHandler.Create())
		}

		charges := gv1.Group("/charges")
		{
			charges.POST("/:paymentOperationId/refund", chargeHandler.Refund())
			charges.GET("/:userID/:serviceType/:serviceID/:paymentOperationID", chargeHandler.GetCharge())
		}

		events := gv1.Group("/events")
		{
			events.POST("/account-update", eventHandler.AccountUpdated())
		}
	}

	service := APIService{
		Engine: router,
		addr:   fmt.Sprintf("%s:%s", host, port),
	}

	return service, nil
}

// Start serves the API
func (api APIService) Start() error {
	return api.Run(api.addr)
}