package parsers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/types"
	"os"
	"reflect"
	"testing"
)

func TestGenerateDateRange(t *testing.T) {
	tests := []struct {
		name     string
		dateFrom string
		dateEnd  string
		result   []string
		wantErr  bool
	}{
		{
			name:     "Successful execution",
			dateFrom: "2024-07-19",
			dateEnd:  "2024-07-21",
			result: []string{
				"2024-07-19",
				"2024-07-20",
				"2024-07-21",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateDateRange(tt.dateFrom, tt.dateEnd)
			if tt.wantErr {
				t.Errorf("GenerateDateRange(%v, %v)", tt.dateFrom, tt.dateEnd)
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equalf(t, tt.result, got, "GenerateDateRange(%v, %v)", tt.dateFrom, tt.dateEnd)
		})
	}
}

func TestFromFiles(t *testing.T) {
	tests := []struct {
		name      string
		dateFrom  string
		dateEnd   string
		want      []types.News
		setup     func(t *testing.T)
		expectErr bool
	}{
		{
			name:     "Valid date range with news",
			dateFrom: "2024-07-19",
			dateEnd:  "2024-07-21",
			want: []types.News{
				{
					Title: "News 1 on 2024-07-19",
				},
				{
					Title: "News 2 on 2024-07-20",
				},
				{
					Title: "News 3 on 2024-07-21",
				},
			},
			setup: func(t *testing.T) {
				newsData := [][]types.News{
					{
						{Title: "News 1 on 2024-07-19"},
					},
					{
						{Title: "News 2 on 2024-07-20"},
					},
					{
						{Title: "News 3 on 2024-07-21"},
					},
				}
				dates := []string{"2024-07-19", "2024-07-20", "2024-07-21"}

				for i, data := range newsData {
					filename := dates[i] + ".json"
					file, err := os.Create(filename)
					assert.Nil(t, err)

					out, err := json.Marshal(data)
					assert.Nil(t, err)

					_, err = file.Write(out)
					assert.Nil(t, err)

					err = file.Close()
					assert.Nil(t, err)
				}
			},
			expectErr: false,
		},
		{
			name:      "Invalid date range",
			dateFrom:  "2024-07-24",
			dateEnd:   "2024-07-23",
			want:      nil,
			setup:     func(t *testing.T) {},
			expectErr: true,
		},
		{
			name:     "Single day range with news",
			dateFrom: "2024-07-23",
			dateEnd:  "2024-07-23",
			want: []types.News{
				{
					Title: "News 1 on 2024-07-23",
				},
				{
					Title: "News 2 on 2024-07-23",
				},
			},
			setup: func(t *testing.T) {
				data := []types.News{
					{
						Title: "News 1 on 2024-07-23",
					},
					{
						Title: "News 2 on 2024-07-23",
					},
				}
				filename := "2024-07-23.json"
				file, err := os.Create(filename)
				assert.Nil(t, err)

				out, err := json.Marshal(data)
				assert.Nil(t, err)

				_, err = file.Write(out)
				assert.Nil(t, err)

				err = file.Close()
				assert.Nil(t, err)
			},
			expectErr: false,
		},
		{
			name:      "No news found",
			dateFrom:  "2024-07-25",
			dateEnd:   "2024-07-25",
			want:      nil,
			setup:     func(t *testing.T) {},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := FromFiles(tt.dateFrom, tt.dateEnd)
			if (err != nil) != tt.expectErr {
				t.Errorf("FromFiles() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
