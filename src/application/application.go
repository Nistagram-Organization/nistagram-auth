package application

import (
	"fmt"
	auth02 "github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/auth_grpc_client"
	"github.com/Nistagram-Organization/nistagram-auth/src/controllers/auth"
	auth2 "github.com/Nistagram-Organization/nistagram-auth/src/services/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	router = gin.Default()
)

const (
	domainKey       = "domain"
	clientIdKey     = "client_id"
	clientSecretKey = "client_secret"
	audienceKey     = "audience"
)

func StartApplication() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	domain := os.Getenv(domainKey)
	clientId := os.Getenv(clientIdKey)
	clientSecret := os.Getenv(clientSecretKey)
	audience := os.Getenv(audienceKey)

	if domain == "" || clientId == "" || clientSecret == "" || audience == "" {
		panic("Environment variables not set properly")
	}

	authController := auth.NewAuthController(
		auth2.NewAuthService(
			auth02.NewAuth0Client(domain, clientId, clientSecret, audience),
			auth_grpc_client.NewAuthGrpcClient(),
		),
	)

	router.POST("/register", authController.Register)

	if port, exists := os.LookupEnv("PORT"); exists {
		router.Run(fmt.Sprintf(":%s", port))
	} else {
		router.Run(":9091")
	}
}
