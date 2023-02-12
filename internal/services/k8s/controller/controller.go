package controller

import (
	"context"
	k8sauth "cron-vault-sync/internal/services/k8s/auth"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"
)

type ObjectsController interface {
	// Secrets
	GetSecret(name string) (*corev1.Secret, error)
	ListSecrets() (*corev1.SecretList, error)

	//ExternalSecrets
	ListExternalSecrets() (*unstructured.UnstructuredList, error)
	CreateExternalSecret(name, keyPath, secretStoreName string) (error)

}

type objectsController struct {
	clientset *kubernetes.Clientset
	dynamicClientSet *dynamic.DynamicClient
	Namespace string
}

func NewObjectsController(namespace string) (ObjectsController, error) {
	clientset, err := k8sauth.NewClient()
	if err != nil {
		return nil, err
	}
	dynamiccs, err := k8sauth.NewDynamicClient()
	if err != nil {
		return nil, err
	}

	
	return &objectsController{
		clientset,
		dynamiccs,
		namespace,
	}, nil


}

func (c *objectsController) GetSecret(name string) (*corev1.Secret, error) {
	return c.clientset.CoreV1().Secrets(c.Namespace).Get(context.Background(),name, metav1.GetOptions{})
}

func (c *objectsController) ListSecrets() (*corev1.SecretList, error) {
	return c.clientset.CoreV1().Secrets(c.Namespace).List(context.Background(),metav1.ListOptions{})
}

func (c *objectsController) CreateExternalSecret(name, keyPath, secretStoreName string) (error) {

	// Define the ExternalSecret CRD
	externalSecret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "external-secrets.io/v1beta1",
			"kind":       "ExternalSecret",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"refreshInterval": "15s",
				"secretStoreRef": map[string]interface{}{
					"name": secretStoreName,
					"kind": "SecretStore",
				},
				"target": map[string]interface{}{
					"name": 	name,
					"namespace": c.Namespace,
				},
				"dataFrom": []map[string]interface{}{
					{
						"extract": map[string]interface{}{
							"key": keyPath,
						},
					},
				},
			},
		},
	}

	gvr := schema.GroupVersionResource{
		Group:    "external-secrets.io",
		Version:  "v1beta1",
		Resource: "externalsecrets",
	}

	
	_, err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).Create(context.Background(), externalSecret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *objectsController) ListExternalSecrets() (*unstructured.UnstructuredList, error) {
	gvr := schema.GroupVersionResource{
		Group:    "external-secrets.io",
		Version:  "v1beta1",
		Resource: "externalsecrets",
	}

	
	result, err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return result, err
}