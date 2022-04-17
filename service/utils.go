package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
)

func GetDeployment(clientset *kubernetes.Clientset, name string, namespace string) *v1.Deployment {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}
	return deployment
}

type Container struct {
	Name string       `json:"name"`
	Env  []v12.EnvVar `json:"env"`
}

type Deployment struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []Container `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

func SetDeploymentEnv(clientset *kubernetes.Clientset, name string, namespace string, containerName string, envList []v12.EnvVar) {
	patch := Deployment{}
	patch.Spec.Template.Spec.Containers = append(patch.Spec.Template.Spec.Containers, Container{Name: containerName, Env: envList})
	log.Info(patch)
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(patchBytes))
	res, err := clientset.AppsV1().Deployments(namespace).Patch(
		context.TODO(), name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(res)
}
