package conformance

import (
	"os"
	"testing"

	"github.com/nauski/cert-manager-webhook-joker/internal/solver"
)

var (
	zone = getEnvOrDefault("TEST_ZONE_NAME", "cert-manager-webhook-joker.test.nauski.fi")
)

func TestRunsSuite(t *testing.T) {
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not specified")
	}

	// TODO: Implement cert-manager conformance tests
	// This requires setting up the full cert-manager test environment
	// For now, test basic solver functionality

	solver := solver.New()
	if solver == nil {
		t.Fatal("Failed to create solver")
	}

	if solver.Name() != "joker" {
		t.Errorf("Expected solver name 'joker', got %s", solver.Name())
	}

	t.Logf("Conformance test placeholder - solver created successfully for domain: %s", zone)
	t.Skip("Full conformance tests require cert-manager test framework setup")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}