package solver

import (
	"testing"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected *Config
	}{
		{
			name: "valid config with secretKeyRef",
			input: `{
				"username": {
					"secretKeyRef": {
						"name": "joker-credentials",
						"key": "username"
					}
				},
				"password": {
					"secretKeyRef": {
						"name": "joker-credentials",
						"key": "password"
					}
				}
			}`,
			wantErr: false,
			expected: &Config{
				Username: &SecretKeySelector{
					SecretKeyRef: &SecretReference{
						Name: "joker-credentials",
						Key:  "username",
					},
				},
				Password: &SecretKeySelector{
					SecretKeyRef: &SecretReference{
						Name: "joker-credentials",
						Key:  "password",
					},
				},
			},
		},
		{
			name:     "empty config",
			input:    "{}",
			wantErr:  false,
			expected: &Config{},
		},
		{
			name:    "invalid json",
			input:   "{invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configJSON := &extapi.JSON{
				Raw: []byte(tt.input),
			}

			got, err := loadConfig(configJSON)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Compare the parsed configuration
			if tt.expected.Username != nil {
				if got.Username == nil || got.Username.SecretKeyRef == nil {
					t.Errorf("loadConfig() username configuration not parsed correctly")
					return
				}
				if got.Username.SecretKeyRef.Name != tt.expected.Username.SecretKeyRef.Name {
					t.Errorf("loadConfig() username name = %v, want %v", got.Username.SecretKeyRef.Name, tt.expected.Username.SecretKeyRef.Name)
				}
				if got.Username.SecretKeyRef.Key != tt.expected.Username.SecretKeyRef.Key {
					t.Errorf("loadConfig() username key = %v, want %v", got.Username.SecretKeyRef.Key, tt.expected.Username.SecretKeyRef.Key)
				}
			}

			if tt.expected.Password != nil {
				if got.Password == nil || got.Password.SecretKeyRef == nil {
					t.Errorf("loadConfig() password configuration not parsed correctly")
					return
				}
				if got.Password.SecretKeyRef.Name != tt.expected.Password.SecretKeyRef.Name {
					t.Errorf("loadConfig() password name = %v, want %v", got.Password.SecretKeyRef.Name, tt.expected.Password.SecretKeyRef.Name)
				}
				if got.Password.SecretKeyRef.Key != tt.expected.Password.SecretKeyRef.Key {
					t.Errorf("loadConfig() password key = %v, want %v", got.Password.SecretKeyRef.Key, tt.expected.Password.SecretKeyRef.Key)
				}
			}
		})
	}
}

func TestLoadConfigNil(t *testing.T) {
	got, err := loadConfig(nil)
	if err != nil {
		t.Errorf("loadConfig(nil) error = %v, want nil", err)
	}
	if got == nil {
		t.Errorf("loadConfig(nil) = nil, want empty config")
	}
}