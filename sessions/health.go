package sessions

import (
	"context"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HealthOrPanic() {
	log.Info("Checking Health ...")
	clientset := GetClientset()
	log.Info("Got Clientset...")
	_, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Got Clientset response...")
}
