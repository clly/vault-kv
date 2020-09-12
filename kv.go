package vaultkv

import "github.com/hashicorp/vault/api"

type KV struct {
	vaultapi *api.Client
}

func New(api *api.Client) *KV {
	return &KV{
		vaultapi: api,
	}
}

func (k *KV) Get(path string) (*api.Secret, error) {
	return k.vaultapi.Logical().Read(path)
}

func (k *KV) List(path string) (*api.Secret, error) {
	return k.vaultapi.Logical().List(path)
}

func (k *KV) Put(path string) (*api.Secret, error) {
	return k.vaultapi.Logical().Write(path, nil)
}

func (k *KV) Delete(path string) (*api.Secret, error) {
	return k.vaultapi.Logical().DeleteWithData(path, nil)
}

func (k *KV)  Raw() *api.Client {
	return k.vaultapi
}
