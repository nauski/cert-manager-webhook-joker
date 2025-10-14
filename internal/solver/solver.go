package solver

import (
	"context"
	"fmt"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/nauski/cert-manager-webhook-joker/internal/joker"
	"github.com/nauski/cert-manager-webhook-joker/internal/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

const (
	GroupName = "acme.joker.com"
)

type Solver struct {
	client kubernetes.Interface
}

func New() *Solver {
	return &Solver{}
}

func (s *Solver) Name() string {
	return "joker"
}

func (s *Solver) Present(ch *v1alpha1.ChallengeRequest) error {
	klog.V(4).Infof("Processing Present request for %s", ch.ResolvedFQDN)

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()
	username, password, err := cfg.loadCredentials(ctx, s.client, ch.ResourceNamespace)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	client := joker.New(joker.Config{
		Username: username,
		Password: password,
	})

	dnsName := util.NormalizeFQDN(ch.ResolvedFQDN)
	zone, label, err := util.ParseChallengeDomain(dnsName)
	if err != nil {
		return fmt.Errorf("failed to parse domain %s: %w", dnsName, err)
	}

	klog.V(4).Infof("Creating TXT record: zone=%s, label=%s, value=%s", zone, label, ch.Key)

	if err := client.CreateTXTRecord(ctx, zone, label, ch.Key); err != nil {
		return fmt.Errorf("failed to create TXT record: %w", err)
	}

	klog.V(4).Infof("Successfully created TXT record for %s", ch.ResolvedFQDN)
	return nil
}

func (s *Solver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	klog.V(4).Infof("Processing CleanUp request for %s", ch.ResolvedFQDN)

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()
	username, password, err := cfg.loadCredentials(ctx, s.client, ch.ResourceNamespace)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	client := joker.New(joker.Config{
		Username: username,
		Password: password,
	})

	dnsName := util.NormalizeFQDN(ch.ResolvedFQDN)
	zone, label, err := util.ParseChallengeDomain(dnsName)
	if err != nil {
		return fmt.Errorf("failed to parse domain %s: %w", dnsName, err)
	}

	klog.V(4).Infof("Deleting TXT record: zone=%s, label=%s", zone, label)

	if err := client.DeleteTXTRecord(ctx, zone, label); err != nil {
		return fmt.Errorf("failed to delete TXT record: %w", err)
	}

	klog.V(4).Infof("Successfully deleted TXT record for %s", ch.ResolvedFQDN)
	return nil
}

func (s *Solver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	klog.V(4).Info("Initializing Joker webhook solver")

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	s.client = cl
	return nil
}
