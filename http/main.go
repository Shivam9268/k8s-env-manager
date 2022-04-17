package http

import (
	"github.com/Shivam9268/k8s-env-manager/http/health"
	"github.com/Shivam9268/k8s-env-manager/http/home"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	healthController *health.Controller
	homeController   *home.Controller
	router           = gin.New()
	server           *http.Server
)

func Initialize() {
	healthController = health.NewController()
	homeController = home.NewController()
	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	setupRoutes()
}

func Shutdown() {
	if err := server.Close(); err != nil {
		log.Fatal("Server Close:", err)
	}
}

func Run() {
	server = &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}

func setupRoutes() {
	router.GET("/health", healthController.HealthHandler)
	router.GET("/", homeController.HomeHandler)
}
