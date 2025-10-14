package solver

import (
	"context"
	"encoding/json"
	"fmt"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	Username *SecretKeySelector `json:"username"`
	Password *SecretKeySelector `json:"password"`
}

type SecretKeySelector struct {
	SecretKeyRef *SecretReference `json:"secretKeyRef"`
}

type SecretReference struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func loadConfig(cfgJSON *extapi.JSON) (*Config, error) {
	cfg := &Config{}
	if cfgJSON == nil {
		return cfg, nil
	}

	if err := json.Unmarshal(cfgJSON.Raw, cfg); err != nil {
		return nil, fmt.Errorf("error decoding solver config: %w", err)
	}

	return cfg, nil
}

func (c *Config) loadCredentials(ctx context.Context, client kubernetes.Interface, namespace string) (username, password string, err error) {
	if c.Username != nil && c.Username.SecretKeyRef != nil {
		username, err = getSecretValue(ctx, client, namespace, c.Username.SecretKeyRef.Name, c.Username.SecretKeyRef.Key)
		if err != nil {
			return "", "", fmt.Errorf("failed to get username from secret: %w", err)
		}
	}

	if c.Password != nil && c.Password.SecretKeyRef != nil {
		password, err = getSecretValue(ctx, client, namespace, c.Password.SecretKeyRef.Name, c.Password.SecretKeyRef.Key)
		if err != nil {
			return "", "", fmt.Errorf("failed to get password from secret: %w", err)
		}
	}

	return username, password, nil
}

func getSecretValue(ctx context.Context, client kubernetes.Interface, namespace, secretName, key string) (string, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretName, err)
	}

	value, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s/%s", key, namespace, secretName)
	}

	return string(value), nil
}
