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
			err := validateFeedSpec(tt.input)
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
		},
		{
			name:    "Invalid URL with no protocol",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL with no protocol",
			wantErr: true,
			url:     string([]byte{0x7f}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlValidator := &urlValidate{
				url: tt.url,
			}
			err := urlValidator.Validate()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
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
			err := validateHotNewsSpec(tt.args.hotNewsSpec)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_baseHandler_HandleNext(t *testing.T) {
	type fields struct {
		next handler
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful validation",
			fields: fields{
				next: &urlValidate{
					url: "http://example.com",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &urlValidate{}
			h.SetNext(tt.fields.next)
			tt.wantErr(t, h.HandleNext(), "HandleNext()")
		})
	}
}

func Test_baseHandler_SetNext(t *testing.T) {
	type fields struct {
		next handler
	}
	type args struct {
		handler handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   baseHandler
	}{
		{
			name: "Successful validation",
			fields: fields{
				next: &urlValidate{},
			},
			want: baseHandler{
				&urlValidate{},
			},
			args: args{
				handler: &urlValidate{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &baseHandler{
				next: tt.fields.next,
			}
			h.SetNext(tt.args.handler)
			assert.Equal(t, tt.want, *h)
		})
	}
}
