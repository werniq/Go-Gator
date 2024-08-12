package v1

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		inputURL  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid name and URL",
			inputName: "ValidName",
			inputURL:  "http://example.com",
			wantErr:   false,
		},
		{
			name:      "Empty name",
			inputName: "",
			inputURL:  "http://example.com",
			wantErr:   true,
			errMsg:    "length of keyword is invalid:  (0)",
		},
		{
			name:      "Name too long",
			inputName: "ThisNameIsWayTooLongToBeValid",
			inputURL:  "http://example.com",
			wantErr:   true,
			errMsg:    "length of keyword is invalid: ThisNameIsWayTooLongToBeValid(27)",
		},
		{
			name:      "Valid name with no HTTP in URL",
			inputName: "ValidName",
			inputURL:  "ftp://example.com",
			wantErr:   true,
			errMsg:    "url must contain http or https",
		},
		{
			name:      "Valid name with no protocol in URL",
			inputName: "ValidName",
			inputURL:  "example.com",
			wantErr:   true,
			errMsg:    "url must contain http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.inputName, tt.inputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, expected error message %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestLengthValidate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		minLen  int
		maxLen  int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid keyword length",
			keyword: "ValidName",
			minLen:  1,
			maxLen:  20,
			wantErr: false,
		},
		{
			name:    "Keyword too short",
			keyword: "",
			minLen:  1,
			maxLen:  20,
			wantErr: true,
			errMsg:  "length of keyword is invalid: (0)",
		},
		{
			name:    "Keyword too long",
			keyword: "ThisNameIsWayTooLongToBeValid",
			minLen:  1,
			maxLen:  20,
			wantErr: true,
			errMsg:  "length of keyword is invalid: ThisNameIsWayTooLongToBeValid(27)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lengthValidator := &LengthValidate{
				keyword:           tt.keyword,
				requiredMinLength: tt.minLen,
				requiredMaxLength: tt.maxLen,
			}
			err := lengthValidator.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LengthValidate.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("LengthValidate.Validate() error = %v, expected error message %v", err.Error(), tt.errMsg)
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
			url:     "ftp://example.com",
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
				t.Errorf("urlValidate.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("urlValidate.Validate() error = %v, expected error message %v", err.Error(), tt.errMsg)
			}
		})
	}
}
