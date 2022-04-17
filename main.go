package main

import (
	"fmt"
	"github.com/Shivam9268/k8s-env-manager/http"
	"github.com/Shivam9268/k8s-env-manager/service"
	"github.com/Shivam9268/k8s-env-manager/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var done = make(chan bool)
var gracefulStop = make(chan os.Signal)

func setup() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error in config file: %s \n", err))
	}
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
	switch viper.GetString("log.level") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		break
	case "INFO":
		log.SetLevel(log.InfoLevel)
		break
	default:
		log.SetLevel(log.ErrorLevel)
		break
	}
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	setupApp()
}

func CleanupOnSignal(cleanup func()) {
	go func() {
		sig := <-gracefulStop
		log.Info(fmt.Sprintf("caught sig: %+v. waiting for goroutines to finish", sig))
		cleanup()
		log.Info("goroutines finished. exiting")
		os.Exit(0)
	}()
}

func setupApp() {
	service.Initialize()
	http.Initialize()
	sessions.HealthOrPanic()
	CleanupOnSignal(cleanup)
}

func cleanup() {
	service.Shutdown()
	http.Shutdown()
	done <- true
}

func run() {
	go http.Run()
	go service.Run()
	<-done
}

func main() {
	setup()
	run()
}
