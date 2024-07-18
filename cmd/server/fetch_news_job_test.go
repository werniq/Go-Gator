package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFetchNewsJob_Run(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_run")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name          string
		input         *FetchNewsJob
		expectDateErr bool
		setup         func()
		finish        func()
	}{
		{
			name: "Default job run",
			input: &FetchNewsJob{
				Filters: types.NewFilteringParams("", time.Now().Format(time.DateOnly), "", ""),
			},
			expectDateErr: false,
			setup:         func() {},
			finish:        func() {},
		},
		{
			name: "Wrong date format",
			input: &FetchNewsJob{
				Filters: types.NewFilteringParams("", time.Now().Format(time.ANSIC), "", ""),
			},
			expectDateErr: true,
			setup:         func() {},
			finish:        func() {},
		},
		{
			name: "Without any sources",
			input: &FetchNewsJob{
				Filters: types.NewFilteringParams("", time.Now().Format(time.DateOnly), "", ""),
			},
			setup: func() {
				supportedSources := parsers.GetAllSources()
				for source, _ := range supportedSources {
					err := parsers.DeleteSource(source)
					if err != nil {
						t.Fatal(err)
					}
				}
			},
			expectDateErr: false,
			finish: func() {
				// revert back all issues caused by setup func, so that next tests will perform correctly
				sourceToEndpoint := map[string]string{
					parsers.WashingtonTimes: "https://www.washingtontimes.com/rss/headlines/news/world",
					parsers.ABC:             "https://abcnews.go.com/abcnews/internationalheadlines",
					parsers.BBC:             "https://feeds.bbci.co.uk/news/rss.xml",
					parsers.UsaToday:        "https://usatoday.com",
				}

				sourceToFormat := map[string]string{
					parsers.UsaToday:        "html",
					parsers.ABC:             "xml",
					parsers.BBC:             "xml",
					parsers.WashingtonTimes: "xml",
				}

				for source, endpoint := range sourceToEndpoint {
					err := parsers.AddNewSource(sourceToFormat[source], source, endpoint)
					if err != nil {
						t.Fatal(err)
					}
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.input.Run()
			if tt.expectDateErr {
				assert.NotNil(t, err)
				return
			} else {
				if err != nil {
					t.Fatalf("Run() returned an error: %v", err)
				}
			}

			CwdPath, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			expectedFilename := filepath.Join(CwdPath, parsers.CmdDir,
				parsers.ParsersDir,
				parsers.DataDir,
				tt.input.Filters.StartingTimestamp+parsers.JsonExtension)

			if _, err := os.Stat(expectedFilename); os.IsNotExist(err) {
				t.Fatalf("expected file %s does not exist", expectedFilename)
			}

			content, err := os.ReadFile(expectedFilename)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			news := []types.News{}
			err = json.Unmarshal(content, &news)
			if err != nil {
				t.Fatalf("failed to unmarshal content: %v", err)
			}

			if news == nil {
				t.Fatalf("news are not expected to be nil")
			}
		})
	}
}
