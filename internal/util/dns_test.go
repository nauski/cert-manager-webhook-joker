package util

import (
	"testing"
)

func TestParseChallengeDomain(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantZone  string
		wantLabel string
		wantErr   bool
	}{
		{
			name:      "subdomain case",
			input:     "int.nauski.fi",
			wantZone:  "nauski.fi",
			wantLabel: "_acme-challenge.int",
			wantErr:   false,
		},
		{
			name:      "root domain case",
			input:     "nauski.fi",
			wantZone:  "nauski.fi",
			wantLabel: "_acme-challenge",
			wantErr:   false,
		},
		{
			name:      "deep subdomain case",
			input:     "test.int.nauski.fi",
			wantZone:  "nauski.fi",
			wantLabel: "_acme-challenge.test.int",
			wantErr:   false,
		},
		{
			name:      "with trailing dot",
			input:     "int.nauski.fi.",
			wantZone:  "nauski.fi",
			wantLabel: "_acme-challenge.int",
			wantErr:   false,
		},
		{
			name:      "invalid domain",
			input:     "invalid",
			wantZone:  "",
			wantLabel: "",
			wantErr:   true,
		},
		{
			name:      "empty domain",
			input:     "",
			wantZone:  "",
			wantLabel: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zone, label, err := ParseChallengeDomain(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChallengeDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if zone != tt.wantZone {
				t.Errorf("ParseChallengeDomain() zone = %v, want %v", zone, tt.wantZone)
			}

			if label != tt.wantLabel {
				t.Errorf("ParseChallengeDomain() label = %v, want %v", label, tt.wantLabel)
			}
		})
	}
}

func TestNormalizeFQDN(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "with trailing dot",
			input: "_acme-challenge.int.nauski.fi.",
			want:  "int.nauski.fi",
		},
		{
			name:  "without trailing dot",
			input: "_acme-challenge.int.nauski.fi",
			want:  "int.nauski.fi",
		},
		{
			name:  "no acme challenge prefix",
			input: "int.nauski.fi",
			want:  "int.nauski.fi",
		},
		{
			name:  "root domain with acme challenge",
			input: "_acme-challenge.nauski.fi",
			want:  "nauski.fi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeFQDN(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeFQDN() = %v, want %v", got, tt.want)
			}
		})
	}
}
