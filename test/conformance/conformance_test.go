package conformance

import (
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/test/acme/dns"
	"github.com/nauski/cert-manager-webhook-joker/internal/solver"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	zone = getEnvOrDefault("TEST_ZONE_NAME", "cert-manager-webhook-joker.test.nauski.fi")
)

func TestRunsSuite(t *testing.T) {
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not specified")
	}

	// The conformance test suite will test:
	// 1. Present() creates the correct TXT record
	// 2. CleanUp() removes the TXT record
	// 3. Handling of multiple simultaneous challenges
	// 4. Proper error handling for invalid credentials
	suite := dns.NewConformanceTestSuite(&dns.ConformanceTestSuiteConfig{
		NewSolver: func() dns.Solver {
			return solver.New()
		},
		ConfigJSONFunc: func() *extapi.JSON {
			return &extapi.JSON{
				Raw: []byte(`{
					"username": {
						"secretKeyRef": {
							"name": "` + os.Getenv("TEST_SECRET_NAME") + `",
							"key": "username"
						}
					},
					"password": {
						"secretKeyRef": {
							"name": "` + os.Getenv("TEST_SECRET_NAME") + `",
							"key": "password"
						}
					}
				}`),
			}
		},
		ManifestPath:        "../../deploy",
		Namespace:           getEnvOrDefault("TEST_NAMESPACE", "cert-manager"),
		DomainName:          zone,
		DNSServerIP:         getEnvOrDefault("TEST_DNS_SERVER", "8.8.8.8:53"),
		UseAuthoritative:    false,
		PropagationTimeout:  300, // 5 minutes for Joker DNS propagation
		PollingInterval:     2,   // Poll every 2 seconds
		StrictMode:          true,
	})

	suite.Run(t)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}