package main

import (
	"cron-vault-sync/internal/services/k8s/controller"
	vaultclient "cron-vault-sync/internal/services/vault"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {

	keyPath := os.Getenv("VAULT_PREFIX_KEY_PATH")
	namespace := os.Getenv("NAMESPACE")
	vclient, err := vaultclient.NewVaultClient()
	if err != nil {
		logrus.Fatal(err)
	}
	ctrl, err := controller.NewObjectsController(namespace)
	if err != nil {
		logrus.Fatal(err)
	}
	secrets, err := vclient.ListSecrets(keyPath)
	if err != nil {
		logrus.Fatal(err)
	}
	vaultCRDSecrets, err := ctrl.ListVaultCRDSecrets()
	if err != nil {
		logrus.Fatal(err)
	}

	secretNames := []string{}
	for _, s := range secrets.Data["keys"].([]interface{}) {
		secretNames = append(secretNames, s.(string))
	}
	crdSecretNames := []string{}
	for _, s := range vaultCRDSecrets.Items {
		crdSecretNames = append(crdSecretNames, s.GetName())
	}

	for _, secretName := range secretNames {
		if !contains(crdSecretNames, secretName) {
			secretKeyPath := keyPath + secretName
			metadata, err := vclient.GetSecretMetadata(secretKeyPath)
			if err != nil {
				logrus.Error(err)
				continue
			}
			customMetadata := make(map[string]interface{})
			if metadata.Data["custom_metadata"] != nil {
				customMetadata = metadata.Data["custom_metadata"].(map[string]interface{})
			}
			if err := ctrl.CreateVaultCRDSecret(secretName, secretKeyPath, customMetadata); err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Created Vault CRD Secret: " + secretName)
			}
		}
	}

	for _, crdSecretName := range crdSecretNames {
		if !contains(secretNames, crdSecretName) {
			if err := ctrl.DeleteVaultCRDSecret(crdSecretName); err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Deleted Vault CRD Secret: " + crdSecretName)
			}
		}
	}
}





func contains(s interface{}, e string) bool {
	switch as := s.(type) {
	case []string:
		for _, a := range as {
			if a == e {
				return true
			}
		}
	case []interface{}:
		for _, a := range as {
			if a.(string) == e {
				return true
			}
		}
	}
	return false
}