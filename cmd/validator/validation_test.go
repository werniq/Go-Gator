package validator

import (
	"testing"
)

func TestValidateDate(t *testing.T) {
	tests := []struct {
		dateStr string
		wantErr bool
	}{
		{"2023-01-01", false},
		{"2023-12-31", false},
		// testing leap year
		{"2023-02-29", true},
		{"", false},
		{"2023/01/01", true},
		{"01-01-2023", true},
	}

	for _, tt := range tests {
		err := ByDate(tt.dateStr)
		if (err != nil) != tt.wantErr {
			t.Errorf("ByDate(%q) = %v, wantErr %v", tt.dateStr, err, tt.wantErr)
		}
	}
}

func TestValidateSources(t *testing.T) {
	tests := []struct {
		sources string
		wantErr bool
	}{
		{"abc,bbc", false},
		{"abc,xyz", true},
		{"usatoday", false},
		{"fakesource", true},
		{"all", false},
	}

	for _, tt := range tests {
		err := BySources(tt.sources)
		if (err != nil) != tt.wantErr {
			t.Errorf("BySources(%v) = %v, wantErr %v", tt.sources, err, tt.wantErr)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice  []string
		item   string
		result bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"abc", "def", "ghi"}, "def", true},
	}

	for _, tt := range tests {
		result := contains(tt.slice, tt.item)
		if result != tt.result {
			t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, result, tt.result)
		}
	}
}
