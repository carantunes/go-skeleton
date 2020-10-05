package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

// Cors returns a new middleware for cors
func Cors(crossOriginDomains string) gin.HandlerFunc {
	return cors.Middleware(cors.Config{
		Origins:        crossOriginDomains,
		Methods:        "GET, PUT, PATCH, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type, Bearer, X-Context-role",
		ExposedHeaders: "Bearer",
		Credentials:    false,
	})
}
