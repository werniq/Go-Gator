package validation

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
