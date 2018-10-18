package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	namespace := "test"
	k8sClient, err := getClient(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	/*
			graph LR
		   svc-service[svc-service] -- pega-http-8080 --> container-name
		   subgraph pod
		   container-name
		   end
	*/

	fmt.Printf("graph LR\n")
	svc, err := getServiceForDeployment("nginx", namespace, k8sClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	pods1, err := getPodsForSvc(svc, namespace, k8sClient)
	_ = pods1
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf(" subgraph containers\n")

	pods, err := getPodsForSvc(svc, namespace, k8sClient)
	_ = pods
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf(" end\n")

}

func getClient(configLocation string) (typev1.CoreV1Interface, error) {
	kubeconfig := filepath.Clean(configLocation)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1(), nil
}

func getServiceForDeployment(deployment string, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.Service, error) {
	listOptions := metav1.ListOptions{}
	svcs, err := k8sClient.Services(namespace).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}
	for _, svc := range svcs.Items {
		if strings.Contains(svc.Name, deployment) {
			// svc-service[svc-service] -- pega-http-8080 --> container-name
			fmt.Fprintf(os.Stdout, " service[%v] -- http-8080 --> ", svc.Name)
			return &svc, nil
		}
	}
	return nil, errors.New("cannot find service for deployment")
}

func getPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := k8sClient.Pods(namespace).List(listOptions)
	for _, pod := range pods.Items {
		fmt.Fprintf(os.Stdout, " %v\n", pod.Name)
	}
	return pods, err
}
