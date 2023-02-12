package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultClient interface {
	ListSecrets(path string) (*vaultapi.Secret, error)
}

type vaultClient struct {
	client *vaultapi.Client
}

func NewVaultClient() (VaultClient, error) {

	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		
		return nil, fmt.Errorf("Error initializing Vault client: %s", err)
	}
	return &vaultClient{client}, nil

}

func (v *vaultClient)ListSecrets(path string) (*vaultapi.Secret, error) {


	secret, err := v.client.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading secrets: %s", err)
	}

	// check secret warnings
	if secret.Warnings != nil {
		return nil, fmt.Errorf("Warnings: %s", secret.Warnings)
	}

	return secret, nil
}