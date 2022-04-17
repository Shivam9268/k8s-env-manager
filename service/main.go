package service

import (
	"github.com/Shivam9268/k8s-env-manager/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	serviceDone      chan bool
	serviceWaitGroup sync.WaitGroup
)

func Initialize() {
	sessions.InitClientset()
	serviceDone = make(chan bool)
}

func Run() {
	serviceWaitGroup.Add(1)
	go startService()
}

func Shutdown() {
	log.Info("not checking configmap anymore ...")
	serviceDone <- true
	serviceWaitGroup.Wait()
}

func startService() {
	defer serviceWaitGroup.Done()
	configMapName := viper.GetString("configmap.name")
	configMapNamespace := viper.GetString("configmap.namespace")
	CheckConfigMap(sessions.GetClientset(), configMapName, configMapNamespace)
	for {
		select {
		case <-time.After(time.Duration(viper.GetInt64("check_interval_seconds")) * time.Second):
			log.Info("checking configmap...")
			CheckConfigMap(sessions.GetClientset(), configMapName, configMapNamespace)
			log.Info("finished checking configmap...")

		case <-serviceDone:
			log.Info("exiting the monitor.")
			return
		}
	}
}
