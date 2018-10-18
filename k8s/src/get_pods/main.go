package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// This program lists the pods in a cluster equivalent to
//
// kubectl get pods
//
func main() {
	// var ns string
	// flag.StringVar(&ns, "namespace", "", "namespace")

	// hardcode the namespace / project for testing
	ns := "test"

	// Bootstrap k8s configuration from local 	Kubernetes config file
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	pods, err := clientset.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}

	services, err := clientset.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}

	// print pods
	fmt.Println("List of pods:")
	for i, pod := range pods.Items {
		// fmt.Printf("[%d] %s\n", i, pod.GetName()+pod.GetLabels())
		fmt.Println(i, pod.GetName(), pod.GetLabels())
	}
	fmt.Print("\n")

	// print services
	fmt.Println("List of services:")
	for i, svc := range services.Items {
		fmt.Println(i, svc.GetName(), svc.GetLabels())
	}
}
