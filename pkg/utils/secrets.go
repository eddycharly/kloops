package utils

import (
	"context"
	"errors"

	"github.com/eddycharly/kloops/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecret(client client.Client, namespace string, secret v1alpha1.Secret) ([]byte, error) {
	if secret.Value != "" {
		return []byte(secret.Value), nil
	}
	nn := types.NamespacedName{
		Namespace: namespace,
		Name:      secret.ValueFrom.SecretKeyRef.Name,
	}
	var s corev1.Secret
	err := client.Get(context.Background(), nn, &s)
	if err != nil {
		return nil, err
	}
	value, ok := s.Data[secret.ValueFrom.SecretKeyRef.Key]
	if !ok {
		return nil, errors.New("key not found in secret")
	}
	return value, nil
}
