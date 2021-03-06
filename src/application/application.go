package application

import (
	auth02 "github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/user_grpc_client"
	"github.com/Nistagram-Organization/nistagram-auth/src/controllers/auth"
	auth2 "github.com/Nistagram-Organization/nistagram-auth/src/services/auth"
	"github.com/Nistagram-Organization/nistagram-auth/src/services/auth_grpc_service"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/prometheus_handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	router        = gin.Default()
	requestsCount = prometheus_handler.GetHttpRequestsCounter()
	requestsSize  = prometheus_handler.GetHttpRequestsSize()
	uniqueUsers   = prometheus_handler.GetUniqueClients()
)

const (
	domainKey       = "domain"
	clientIdKey     = "client_id"
	clientSecretKey = "client_secret"
	audienceKey     = "audience"
	dockerKey       = "docker"
)

func configureCORS() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))
}

func registerPrometheusMiddleware() {
	prometheus.Register(requestsCount)
	prometheus.Register(requestsSize)
	prometheus.Register(uniqueUsers)

	router.Use(prometheus_handler.PrometheusMiddleware(requestsCount, requestsSize, uniqueUsers))
}

func StartApplication() {
	configureCORS()
	registerPrometheusMiddleware()

	domain := os.Getenv(domainKey)
	clientId := os.Getenv(clientIdKey)
	clientSecret := os.Getenv(clientSecretKey)
	audience := os.Getenv(audienceKey)
	var docker bool
	if os.Getenv(dockerKey) == "" {
		docker = false
	} else {
		docker = true
	}

	if domain == "" || clientId == "" || clientSecret == "" || audience == "" {
		panic("Environment variables not set properly")
	}

	port := ":9091"
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	auth0Client := auth02.NewAuth0Client(domain, clientId, clientSecret, audience)
	userGrpcClient := user_grpc_client.NewUserGrpcClient(docker)
	authService := auth2.NewAuthService(
		auth0Client,
		userGrpcClient,
	)
	authController := auth.NewAuthController(authService)

	grpcS := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcS,
		auth_grpc_service.NewAuthGrpcService(auth0Client),
	)

	router.POST("/register", authController.Register)

	router.GET("/metrics", prometheus_handler.PrometheusGinHandler())

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
