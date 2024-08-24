package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   FeedSpec
		wantErr bool
	}{
		{
			name: "Valid name and URL",
			input: FeedSpec{
				Name: "ValidName",
				Link: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid name with no HTTP in URL",
			input: FeedSpec{
				Name: "ValidName",
				Link: "bbc:/example.com",
			},
			wantErr: true,
		},
		{
			name: "Valid name with no protocol in URL",
			input: FeedSpec{
				Name: "ValidName",
				Link: "ftp:/examplecom",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeeds(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUrlValidate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid URL with http",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "Valid URL with https",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "Invalid URL with no http/https",
			url:     "ftp:////example.com",
			wantErr: true,
			errMsg:  "url must contain http or https",
		},
		{
			name:    "Invalid URL with no protocol",
			url:     "example.com",
			wantErr: true,
			errMsg:  "url must contain http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlValidator := &urlValidate{
				url: tt.url,
			}
			err := urlValidator.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("urlValidate.validateFeeds() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("urlValidate.validateFeeds() error = %v, expected error message %v", err.Error(), tt.errMsg)
			}
		})
	}
}
