package dependencies

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/carantunes/go-skeleton/mongodb"
	"github.com/carantunes/go-skeleton/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/uniplaces/stripe-gateway/domain"
	"github.com/uniplaces/stripe-gateway/domain/account"
	"github.com/uniplaces/stripe-gateway/domain/charge"
	"github.com/uniplaces/stripe-gateway/domain/customer"
	"github.com/uniplaces/stripe-gateway/domain/event"
	"github.com/uniplaces/stripe-gateway/domain/source"
	"github.com/uniplaces/stripe-gateway/domain/transfergroup"
	"github.com/uniplaces/stripe-gateway/dynamodb"
	"github.com/uniplaces/stripe-gateway/gin"
	"github.com/uniplaces/stripe-gateway/stripe"
)

// Dependencies holds exported dependencies
type Dependencies struct {
	API gin.APIService
	DB  dynamodb.DB
}

// New resolves and returns initialized Dependencies
func New() (Dependencies, error) {
	cfg, err := viper.New(os.Getenv)
	if err != nil {
		return Dependencies{}, err
	}

	db, err := newDBService(cfg)
	if err != nil {
		return Dependencies{}, err
	}

	api, err := gin.New(
		cfg.GetString("app.host"),
		cfg.GetString("app.port"),
		cfg.GetString("app.cors_domain"),
		cfg.GetString("app.mode"),
		newInteractorAggregate(cfg, db),
	)
	if err != nil {
		return Dependencies{}, err
	}

	return Dependencies{
		API: api,
		DB:  db,
	}, nil
}

func newInteractorAggregate(cfg viper.ConfigService, db dynamodb.Service) domain.InteractorAggregate {
	counterStorage := mongodb.NewCounterStorageService(db)

	// Stripe provider services
	paymentServiceProvider := stripe.NewPaymentProviderService(
		cfg.GetString("STRIPE_API_KEY"),
		cfg.GetString("GOENV"),
		cfg.GetString("STRIPE_LOCAL_URL"),
	)
	customerProvider := stripe.NewCustomerService(paymentServiceProvider)
	sourceProvider := stripe.NewSourceService(paymentServiceProvider)
	chargeProvider := stripe.NewChargeService(paymentServiceProvider)
	refundProvider := stripe.NewRefundService(paymentServiceProvider)
	transferProvider := stripe.NewTransferService(paymentServiceProvider)
	payoutProvider := stripe.NewPayoutService(paymentServiceProvider)
	eventProvider := stripe.NewEventService(cfg.GetString("STRIPE_WEBHOOK_ACCOUNT_UPDATE_SECRET"))

	// Interactor services
	transferGroupInteractor := transfergroup.NewInteractor(transferGroupStorage)

	return domain.InteractorAggregate{
		Customer: customer.NewInteractor(
			customerProvider,
			chargeProvider,
			customerStorage,
			transferGroupInteractor,
			chargeStorage,
		),
		Source: source.NewInteractor(sourceProvider),
		Charge: charge.NewInteractor(refundProvider, chargeProvider, chargeStorage),
		Account: account.NewInteractor(
			transferProvider,
			accountStorage,
			transferGroupInteractor,
			transferStorage,
			payoutProvider,
			payoutStorage,
		),
		Event: event.NewInteractor(accountStorage, eventProvider),
	}
}

func newDBService(cfg viper.ConfigService) (mongodb.DB, error) {
	port := cfg.GetString("db.port")
	endpoint := cfg.GetString("db.endpoint")
	timeout := cfg.GetDuration("db.timeout")

	ctx, _ := context.WithTimeout(context.Background(), timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", endpoint, port)))

	// DynamoDB services
	db, err := mongodb.New(client)
	if err != nil {
		return mongodb.DB{}, err
	}

	return db, nil
}
