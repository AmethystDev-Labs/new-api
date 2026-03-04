package common

import "testing"

func TestValidateAMFSBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{
			name:    "valid default domain",
			rawURL:  "https://amfs.amethyst.ltd",
			wantErr: false,
		},
		{
			name:    "valid subdomain",
			rawURL:  "https://sub.amfs.amethyst.ltd",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			rawURL:  "http://amfs.amethyst.ltd",
			wantErr: true,
		},
		{
			name:    "invalid untrusted domain",
			rawURL:  "https://evil.example.com",
			wantErr: true,
		},
		{
			name:    "invalid with query",
			rawURL:  "https://amfs.amethyst.ltd?x=1",
			wantErr: true,
		},
		{
			name:    "invalid with path",
			rawURL:  "https://amfs.amethyst.ltd/api",
			wantErr: true,
		},
		{
			name:    "invalid with custom port",
			rawURL:  "https://amfs.amethyst.ltd:8443",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateAMFSBaseURL(tt.rawURL)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})
	}
}
