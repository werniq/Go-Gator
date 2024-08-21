package controller

import (
	"testing"
)

// TestServerErr tests the Error method of the serverErr struct.
func TestServerErr(t *testing.T) {
	testCases := []struct {
		name        string
		errorMsg    string
		expectedErr string
	}{
		{
			name:        "Non-empty error message",
			errorMsg:    "Something went wrong",
			expectedErr: "Something went wrong",
		},
		{
			name:        "Empty error message",
			errorMsg:    "",
			expectedErr: "",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &serverErr{
				ErrorMsg: tc.errorMsg,
			}

			actualErr := err.Error()

			if actualErr != tc.expectedErr {
				t.Errorf("expected error message: %s, got: %s", tc.expectedErr, actualErr)
			}
		})
	}
}
