package validator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"newsaggr/cmd/parsers"
	"testing"
)

func TestDateRangeHandler_Handle(t *testing.T) {
	tests := []struct {
		name      string
		dateFrom  string
		dateEnd   string
		expectErr error
	}{
		{"ValidDateRange", "2024-05-01", "2024-05-15", nil},
		{"InvalidDateRange", "2024-05-16", "2024-05-15", errors.New(ErrDateFromAfter)},
		{"EmptyDates", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &DateRangeHandler{
				dateFrom: tt.dateFrom,
				dateEnd:  tt.dateEnd,
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, "/?date-from="+tt.dateFrom+"&date-end="+tt.dateEnd, nil)

			err := handler.Handle()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestDateValidationHandler_Handle(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		expectErr error
	}{
		{"ValidDates", "2024-05-01", nil},
		{"InvalidDateFrom", "2024-15-01", errors.New(ErrFailedDateValidation)},
		{"InvalidDateEnd", "2024-05-51", errors.New(ErrFailedDateValidation)},
		{"EmptyDates", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &DateValidationHandler{
				date: tt.date,
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, "/?date-from="+tt.date, nil)

			err := handler.Handle()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestSourceValidationHandler_Handle(t *testing.T) {
	tests := []struct {
		name      string
		sources   string
		expectErr error
	}{
		{"ValidSources", "abc,bbc", nil},
		{"InvalidSource", "abc,xyz", errors.New(
			fmt.Sprintf("%v%s%v", ErrFailedSourceValidation, "unsupported source: xyz. Supported sources are: ",
				parsers.GetAllSources()))},

		{"EmptySources", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &SourceValidationHandler{
				sources: tt.sources,
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, "/?sources="+tt.sources, nil)

			err := handler.Handle()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

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
		{"firstS,secondS", true},
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
