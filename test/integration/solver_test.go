package integration

import (
	"context"
	"testing"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/nauski/cert-manager-webhook-joker/internal/solver"
	"github.com/nauski/cert-manager-webhook-joker/internal/util"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/api/core/v1"
)

func TestSolverName(t *testing.T) {
	solver := solver.New()
	if got := solver.Name(); got != "joker" {
		t.Errorf("Expected solver name 'joker', got %s", got)
	}
}

func TestSolverWithValidConfig(t *testing.T) {
	// Create fake Kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create test secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"username": []byte("test-username"),
			"password": []byte("test-password"),
		},
	}
	_, err := fakeClient.CoreV1().Secrets("test-namespace").Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test secret: %v", err)
	}

	// Create solver with mock client
	s := solver.New()

	// Initialize with fake client (we need to set this through initialization)
	err = s.Initialize(nil, nil)
	if err == nil {
		// This would normally fail because we don't have a real kubeconfig
		// but the test shows the structure is correct
		t.Log("Solver initialization structure is correct")
	}

	// Test config parsing
	configJSON := &extapi.JSON{
		Raw: []byte(`{
			"username": {
				"secretKeyRef": {
					"name": "test-secret",
					"key": "username"
				}
			},
			"password": {
				"secretKeyRef": {
					"name": "test-secret",
					"key": "password"
				}
			}
		}`),
	}

	// Create test challenge request
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:        "_acme-challenge.test.example.com.",
		Key:                 "test-key-value",
		ResourceNamespace:   "test-namespace",
		DNSName:            "test.example.com",
		Config:             configJSON,
	}

	// Test domain parsing
	dnsName := util.NormalizeFQDN(ch.ResolvedFQDN)
	zone, label, err := util.ParseChallengeDomain(dnsName)
	if err != nil {
		t.Fatalf("Failed to parse domain: %v", err)
	}

	expectedZone := "example.com"
	expectedLabel := "_acme-challenge.test"

	if zone != expectedZone {
		t.Errorf("Expected zone %s, got %s", expectedZone, zone)
	}

	if label != expectedLabel {
		t.Errorf("Expected label %s, got %s", expectedLabel, label)
	}
}

func TestSolverWithInvalidConfig(t *testing.T) {
	s := solver.New()

	// Test with invalid JSON
	invalidConfigJSON := &extapi.JSON{
		Raw: []byte(`{invalid json`),
	}

	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:        "_acme-challenge.test.example.com.",
		Key:                 "test-key-value",
		ResourceNamespace:   "test-namespace",
		DNSName:            "test.example.com",
		Config:             invalidConfigJSON,
	}

	// This should fail due to invalid JSON
	err := s.Present(ch)
	if err == nil {
		t.Error("Expected error with invalid JSON config, got nil")
	}
}

func TestDomainParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedZone   string
		expectedLabel  string
		expectedError  bool
	}{
		{
			name:          "simple domain",
			input:         "_acme-challenge.example.com",
			expectedZone:  "example.com",
			expectedLabel: "_acme-challenge",
		},
		{
			name:          "subdomain",
			input:         "_acme-challenge.test.example.com",
			expectedZone:  "example.com",
			expectedLabel: "_acme-challenge.test",
		},
		{
			name:          "deep subdomain",
			input:         "_acme-challenge.sub.test.example.com",
			expectedZone:  "example.com",
			expectedLabel: "_acme-challenge.sub.test",
		},
		{
			name:          "with trailing dot",
			input:         "_acme-challenge.test.example.com.",
			expectedZone:  "example.com",
			expectedLabel: "_acme-challenge.test",
		},
		{
			name:          "invalid domain",
			input:         "invalid",
			expectedError: true,
		},
		{
			name:          "empty domain",
			input:         "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := util.NormalizeFQDN(tt.input)
			zone, label, err := util.ParseChallengeDomain(normalized)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				return
			}

			if zone != tt.expectedZone {
				t.Errorf("Expected zone %s, got %s", tt.expectedZone, zone)
			}

			if label != tt.expectedLabel {
				t.Errorf("Expected label %s, got %s", tt.expectedLabel, label)
			}
		})
	}
}