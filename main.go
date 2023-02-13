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
		// if secret only exists in vault, delete it from k8s
		for _, vs := range vaultCRDSecrets {
			if !contains(secret.Data["keys"].([]interface{}), vs) {
				err = ctrl.DeleteVaultCRDSecret(vs)
				if err != nil {
					logrus.Error(err)
				} else {
					logrus.Info("Deleted Vault CRD Secret: " + vs)
				}
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