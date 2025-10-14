package util

import (
	"fmt"
	"strings"
)

func ParseChallengeDomain(dnsName string) (zone, label string, err error) {
	if dnsName == "" {
		return "", "", fmt.Errorf("dns name cannot be empty")
	}

	dnsName = strings.TrimSuffix(dnsName, ".")

	parts := strings.Split(dnsName, ".")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid domain name: %s", dnsName)
	}

	zone = strings.Join(parts[len(parts)-2:], ".")

	subdomain := strings.Join(parts[:len(parts)-2], ".")
	if subdomain == "" {
		label = "_acme-challenge"
	} else {
		label = "_acme-challenge." + subdomain
	}

	return zone, label, nil
}

func NormalizeFQDN(fqdn string) string {
	fqdn = strings.TrimSuffix(fqdn, ".")

	if strings.HasPrefix(fqdn, "_acme-challenge.") {
		return strings.TrimPrefix(fqdn, "_acme-challenge.")
	}

	return fqdn
}
