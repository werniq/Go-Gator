package news_fetcher

import (
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"gogator/cmd/server/handlers"
	"gogator/cmd/types"
	"os"
	"testing"
	"time"
)

func TestRunJob(t *testing.T) {
	storagePath := parsers.StoragePath
	tests := []struct {
		name      string
		job       *NewsFetchingJob
		expectErr bool
		setup     func()
		finish    func()
	}{
		{
			name: "Successful job execution",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			expectErr: false,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "Invalid storage path",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			expectErr: true,
			setup: func() {
				parsers.StoragePath = "/invalid_path/"
			},
			finish: func() {
				parsers.StoragePath = storagePath
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := RunJob()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, time.Now().Format("2006-01-02"), handlers.LastFetchedFileDate)
			}

			tt.finish()
		})
	}
}

func TestFetchingJob_Execute(t *testing.T) {
	storagePath := parsers.StoragePath
	tempDir, err := os.MkdirTemp("", "test_execute")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("failed to remove temp directory: %v", err)
		}
	}(tempDir)

	testCases := []struct {
		name      string
		job       *NewsFetchingJob
		expectErr bool
		setup     func()
		finish    func()
	}{
		{
			name: "Default job run",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			expectErr: false,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "Invalid date format",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format(time.ANSIC), "", ""),
			},
			expectErr: true,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "Empty filter parameters",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", "", "", ""),
			},
			expectErr: false,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "Invalid storage path",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			expectErr: true,
			setup: func() {
				parsers.StoragePath = "/invalid_path/"
			},
			finish: func() {
				parsers.StoragePath = storagePath
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer tc.finish()

			err := tc.job.Execute()
			if tc.expectErr {
				assert.NotNil(t, err)
				return
			} else {
				if err != nil {
					t.Fatalf("Execute() returned an error: %v", err)
				}
			}
		})
	}
}
