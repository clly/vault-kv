package vaultkv

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/vault/api"
)

const internalUIPath = "/v1/sys/internal/ui/mounts"

type KV struct {
	vaultapi *api.Client
}

func New(client *api.Client) *KV {
	return &KV{
		vaultapi: client,
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

func (k *KV) Raw() *api.Client {
	return k.vaultapi
}

func isVersionedKV(path string, client api.Client) (bool, error) {
	req := client.NewRequest("GET", internalUIPath)
	resp, err := client.RawRequest(req)
	if err != nil {
		return false, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return false, err
	}
	secretMounts, ok := secret.Data["secret"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("failed to retrieve secret mounts")
	}
	mountPath := getMountFromPath(path)
	rawMount, ok := secretMounts[mountPath]
	if !ok {
		return false, nil
	}

	mount := &api.MountOutput{}
	err = mapstructure.Decode(rawMount, api.MountOutput{})
	if err != nil {
		return false, fmt.Errorf("failed to decode mount configuration %w", err)
	}
	v, ok := mount.Options["version"]
	if !ok {
		return false, nil
	}
	if v == "2" {
		return true, nil
	}
	return false, nil
}

func getMountFromPath(path string) string {
	idx := strings.Index(path, "/")
	if idx == -1 {
		return path
	}
	return path[:idx]
}
