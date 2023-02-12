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
		logrus.New().Fatal(err)
	}

	ctrl, err := controller.NewObjectsController(namespace)
	if err != nil {
	 	logrus.Fatal(err)
	}

	secret, err := vclient.ListSecrets(keyPath)
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := ctrl.ListVaultCRDSecrets()
	if err != nil {
		logrus.Fatal(err)
	}


	vaultCRDSecrets := []string{}
	for _, s := range result.Items {
		vaultCRDSecrets = append(vaultCRDSecrets, s.GetName())
	}

	for _, s := range secret.Data["keys"].([]interface{}) {
		secretName := s.(string)
		if !contains(vaultCRDSecrets, secretName) {
			secretKeyPath := keyPath + s.(string)
			secretMetadata, err := vclient.GetSecretMetadata(secretKeyPath)
			if err != nil {
				logrus.Error(err)
			}
			
			var customMetadata map[string]interface{}
			if secretMetadata.Data["custom_metadata"] != nil {
				customMetadata = secretMetadata.Data["custom_metadata"].(map[string]interface{})
			}

			err = ctrl.CreateVaultCRDSecret(secretName, secretKeyPath, customMetadata)
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Created Vault CRD Secret: " + secretName)
			}
		}

		
	}
}

func contains[T []string](s T, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}