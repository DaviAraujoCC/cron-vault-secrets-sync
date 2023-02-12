package auth

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	kubeconfig = clientcmd.RecommendedHomeFile
	configk8s  *rest.Config
	err        error
)

func NewClient() (*kubernetes.Clientset, error) {

	configk8s, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		configk8s, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(configk8s)
	if err != nil {
		return nil, err
	}



	return clientset, nil
}

func NewDynamicClient() (*dynamic.DynamicClient, error) {

	configk8s, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		configk8s, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	
	if err != nil {
		return nil, err
	}

	dynamiccs, err := dynamic.NewForConfig(configk8s)
	if err != nil {
		return nil, err
	}



	return dynamiccs, nil
}