package news_fetcher

import (
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
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
		args      string
		expectErr bool
		setup     func()
		finish    func()
	}{
		{
			name: "Successful job execution",
			args: storagePath,
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
			args:      "\\invalid-path\\invalid-dir\\",
			expectErr: true,
			setup: func() {
				parsers.StoragePath = ""
			},
			finish: func() {
				parsers.StoragePath = storagePath
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := RunJob(tt.args)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
		args      string
		expectErr bool
	}{
		{
			name: "Default job run",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			args:      storagePath,
			expectErr: false,
		},
		{
			name: "Invalid date format",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format(time.ANSIC), "", ""),
			},
			args:      storagePath,
			expectErr: true,
		},
		{
			name: "Invalid storage path",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format("2006-01-02"), "", ""),
			},
			args:      "\\invalid-path\\invalid-dir\\",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.job.storagePath = tc.args

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
