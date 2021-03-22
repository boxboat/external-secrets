package aws

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	esv1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	"github.com/external-secrets/external-secrets/pkg/provider"
	"github.com/external-secrets/external-secrets/pkg/provider/aws/parameterstore"
	"github.com/external-secrets/external-secrets/pkg/provider/aws/secretsmanager"
	awssess "github.com/external-secrets/external-secrets/pkg/provider/aws/session"
	"github.com/external-secrets/external-secrets/pkg/provider/schema"
)

// Provider satisfies the provider interface.
type Provider struct{}

// NewClient constructs a new secrets client based on the provided store.
func (p *Provider) NewClient(ctx context.Context, store esv1alpha1.GenericStore, kube client.Client, namespace string) (provider.SecretsClient, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	spec := store.GetSpec()
	if spec == nil {
		return nil, fmt.Errorf("store is missing spec")
	}
	if spec.Provider == nil {
		return nil, fmt.Errorf("storeSpec is missing provider")
	}
	if spec.Provider.AWS == nil {
		return nil, fmt.Errorf("storeSpec is missing aws spec")
	}
	switch spec.Provider.AWS.Service {
	case esv1alpha1.AWSServiceSecretsManager:
		return secretsmanager.New(ctx, store, kube, namespace, awssess.DefaultSTSProvider)
	case esv1alpha1.AWSServiceParameterStore:
		return parameterstore.New(ctx, store, kube, namespace, awssess.DefaultSTSProvider)
	}
	return nil, fmt.Errorf("unknown AWS Provider Service: %s", spec.Provider.AWS.Service)
}

func init() {
	schema.Register(&Provider{}, &esv1alpha1.SecretStoreProvider{
		AWS: &esv1alpha1.AWSProvider{},
	})
}