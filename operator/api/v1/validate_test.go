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

func Test_dateValidate_Validate(t *testing.T) {
	type fields struct {
		baseHandler baseHandler
		dateStart   string
		dateEnd     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Successful validation",
			fields: fields{
				baseHandler: baseHandler{},
				dateStart:   "2024-06-01",
				dateEnd:     "2024-06-02",
			},
		},
		{
			name: "Invalid date range",
			fields: fields{
				baseHandler: baseHandler{},
				dateStart:   "2024-06-02",
				dateEnd:     "2024-06-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid date start format",
			fields: fields{
				baseHandler: baseHandler{},
				dateStart:   "AAAAAAAAAA",
				dateEnd:     "2024-06-02",
			},
			wantErr: true,
		},
		{
			name: "Invalid date end format",
			fields: fields{
				baseHandler: baseHandler{},
				dateStart:   "2024-06-01",
				dateEnd:     "AAAAAAAAAA",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dateValidate{
				baseHandler: tt.fields.baseHandler,
				dateStart:   tt.fields.dateStart,
				dateEnd:     tt.fields.dateEnd,
			}
			got := d.Validate()
			if tt.wantErr {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_validateHotNews(t *testing.T) {
	type args struct {
		hotNewsSpec HotNewsSpec
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful validation",
			args: args{
				hotNewsSpec: HotNewsSpec{
					DateStart: "2024-06-01",
					DateEnd:   "2024-06-02",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid date range",
			args: args{
				hotNewsSpec: HotNewsSpec{
					DateStart: "2024-06-02",
					DateEnd:   "2024-06-01",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHotNews(tt.args.hotNewsSpec)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
