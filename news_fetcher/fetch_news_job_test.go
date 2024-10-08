package main

import (
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"os"
	"path/filepath"
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
		time.Sleep(time.Second * 5)
		err := os.Remove(filepath.Join(path, time.Now().Format(time.DateOnly)+".json"))
		if err != nil {
			t.Fatalf("failed to remove data file: %v", err)
		}
		err = os.RemoveAll(path)
		if err != nil {
			t.Fatalf("failed to remove temp directory: %v", err)
		}
	}(tempDir)

	testCases := []struct {
		name      string
		job       *NewsFetchingJob
		args      string
		expectErr bool
		setup     func()
		finish    func()
	}{
		{
			name: "Default job run",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format(time.DateOnly), "", ""),
			},
			args:      storagePath,
			expectErr: false,
			setup:     func() {},
			finish: func() {
				err := os.Remove(time.Now().Format(time.DateOnly) + ".json")
				assert.Nil(t, err)
			},
		},
		{
			name: "Invalid storage path",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format(time.DateOnly), "", ""),
			},
			args:      "\\invalid-path\\invalid-dir\\",
			expectErr: true,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "File creation error",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", string([]byte{0x00, 0x3C, 0x3E, 0x7C}), "", ""),
			},
			args:      tempDir,
			expectErr: true,
			setup:     func() {},
			finish:    func() {},
		},
		{
			name: "Parse by source error",
			job: &NewsFetchingJob{
				params: types.NewFilteringParams("", time.Now().Format(time.DateOnly), "", ""),
			},
			args:      tempDir,
			expectErr: true,
			setup: func() {
				err := parsers.AddNewSource("xml", "nonexistent", "nonexistent")
				assert.Nil(t, err)
			},
			finish: func() {
				err := parsers.DeleteSource("nonexistent")
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer tc.finish()
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
