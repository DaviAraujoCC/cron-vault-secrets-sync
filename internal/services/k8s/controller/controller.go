package controller

import (
	"context"
	k8sauth "cron-vault-sync/internal/services/k8s/auth"
	"strings"

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

	//Vault CRD Secrets
	GetVaultCRDSecret(name string) (*unstructured.Unstructured, error)
	ListVaultCRDSecrets() (*unstructured.UnstructuredList, error)
	CreateVaultCRDSecret(name, secretKeyPath string, customMetadata map[string]interface{}) (error)
	UpdateVaultCRDSecret(name string, vaultCRDSecret *unstructured.Unstructured) (error)
	DeleteVaultCRDSecret(name string) (error)

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

func (c *objectsController) GetVaultCRDSecret(name string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "koudingspawn.de",
		Version:  "v1",
		Resource: "vault",
	}

	return c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *objectsController) CreateVaultCRDSecret(name, secretKeyPath string, customMetadata map[string]interface{}) (error) {

	secretKeyPath = strings.ReplaceAll(secretKeyPath, "metadata/", "")

	vaultcrdsecret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "koudingspawn.de/v1",
			"kind":       "Vault",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"path": secretKeyPath,
				"type": "KEYVALUEV2",
			},
		},
	}
	
	if customMetadata != nil {
		if v, ok := customMetadata["app-owner"]; ok {
			vaultcrdsecret.Object["spec"].(map[string]interface{})["changeAdjustmentCallback"] = map[string]interface{}{
				"type": "deployment",
				"name": v,
			}
		}
	}
	

	gvr := schema.GroupVersionResource{
		Group:    "koudingspawn.de",
		Version:  "v1",
		Resource: "vault",
	}

	
	_, err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).Create(context.Background(), vaultcrdsecret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *objectsController) ListVaultCRDSecrets() (*unstructured.UnstructuredList, error) {
	gvr := schema.GroupVersionResource{
		Group:    "koudingspawn.de",
		Version:  "v1",
		Resource: "vault",
	}

	
	result, err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *objectsController) UpdateVaultCRDSecret(name string, vaultCRDSecret *unstructured.Unstructured) (error) {
	gvr := schema.GroupVersionResource{
		Group:    "koudingspawn.de",
		Version:  "v1",
		Resource: "vault",
	}

	_, err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).Update(context.Background(), vaultCRDSecret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *objectsController) DeleteVaultCRDSecret(name string) (error) {

	gvr := schema.GroupVersionResource{
		Group:    "koudingspawn.de",
		Version:  "v1",
		Resource: "vault",
	}

	err := c.dynamicClientSet.Resource(gvr).Namespace(c.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}