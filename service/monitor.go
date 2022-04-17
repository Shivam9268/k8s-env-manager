package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"reflect"
)

func CheckConfigMap(clientset *kubernetes.Clientset, name string, namespace string) {
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln(configMap.Data)
	for deploymentName, env := range configMap.Data {
		var parsedEnv []map[string]string
		err = json.Unmarshal([]byte(env), &parsedEnv)
		if err != nil {
			log.Fatal(err)
		}
		CheckDeploymentEnv(clientset, deploymentName, namespace, parsedEnv)
	}
}

func CheckDeploymentEnv(clientset *kubernetes.Clientset, name string, namespace string, configMapEnv []map[string]string) {
	deployment := GetDeployment(clientset, name, namespace)
	if deployment == nil {
		log.Infof("Deployment %s not found, hence skipping", name)
		return
	}
	mergedEnv := map[string]v12.EnvVar{}
	currentEnv := map[string]v12.EnvVar{}
	for _, envVar := range deployment.Spec.Template.Spec.Containers[0].Env {
		mergedEnv[envVar.Name] = envVar
		currentEnv[envVar.Name] = envVar
	}
	for _, envVar := range configMapEnv {
		if _, ok := envVar["name"]; !ok {
			log.Fatal("Env does not contain key name")
		}
		if _, ok := envVar["value"]; !ok {
			log.Fatal("Env does not contain key value")
		}
		mergedEnv[envVar["name"]] = v12.EnvVar{Name: envVar["name"], Value: envVar["value"]}
	}
	if !reflect.DeepEqual(mergedEnv, currentEnv) {
		log.WithFields(log.Fields{
			"env":                  configMapEnv,
			"deployment_name":      name,
			"deployment_namespace": namespace,
		}).Info("Updating deployment environment variables")
		var envList []v12.EnvVar
		for _, value := range mergedEnv {
			envList = append(envList, value)
		}
		SetDeploymentEnv(clientset, name, namespace, deployment.Spec.Template.Spec.Containers[0].Name, envList)
	} else {
		log.Infof("Skipping deployment %s as env match", name)
	}
}
