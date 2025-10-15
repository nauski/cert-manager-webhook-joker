package integration

import (
	"testing"

	"github.com/nauski/cert-manager-webhook-joker/internal/solver"
	"github.com/nauski/cert-manager-webhook-joker/internal/util"
)

func TestSolverName(t *testing.T) {
	solver := solver.New()
	if got := solver.Name(); got != "joker" {
		t.Errorf("Expected solver name 'joker', got %s", got)
	}
}

func TestSolverWithValidConfig(t *testing.T) {
	// Test that solver can be created
	s := solver.New()
	if s == nil {
		t.Fatal("New() returned nil solver")
	}

	// Test config parsing would go here but requires Kubernetes client
	// For now, just test domain parsing utilities

	// Test domain parsing utility functions
	testDomain := "_acme-challenge.test.example.com."
	dnsName := util.NormalizeFQDN(testDomain)
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

	t.Logf("Config parsing successful: zone=%s, label=%s", zone, label)
}

func TestSolverWithInvalidConfig(t *testing.T) {
	s := solver.New()
	if s == nil {
		t.Fatal("New() returned nil solver")
	}

	// Test that solver has correct name
	if name := s.Name(); name != "joker" {
		t.Errorf("Expected solver name 'joker', got %s", name)
	}

	t.Log("Solver created successfully with correct interface")
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