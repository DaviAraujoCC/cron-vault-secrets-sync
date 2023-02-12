package main

import (
	"cron-vault-sync/internal/services/k8s/controller"
	vaultclient "cron-vault-sync/internal/services/vault"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {

	path := os.Getenv("VAULT_PREFIX_KEY_PATH")
	secretStoreRef := os.Getenv("SECRET_STORE_REF")
	namespace := os.Getenv("NAMESPACE")

	vclient, err := vaultclient.NewVaultClient()
	if err != nil {
		logrus.New().Fatal(err)
	}

	ctrl, err := controller.NewObjectsController(namespace)
	if err != nil {
	 	logrus.Fatal(err)
	}

	secret, err := vclient.ListSecrets(path)
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := ctrl.ListExternalSecrets()
	if err != nil {
		logrus.Fatal(err)
	}


	externalSecrets := []string{}
	for _, s := range result.Items {
		externalSecrets = append(externalSecrets, s.GetName())
	}

	for _, s := range secret.Data["keys"].([]interface{}) {
		secretName := s.(string)
		if !contains(externalSecrets, secretName) {
			keyPath := strings.ReplaceAll(path + s.(string), "metadata/", "")
			err = ctrl.CreateExternalSecret(secretName, keyPath, secretStoreRef)
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Created ExternalSecret: " + secretName)
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