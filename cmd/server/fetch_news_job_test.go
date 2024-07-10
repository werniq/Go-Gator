package server

import (
	"encoding/json"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
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

	d := time.Now().Format(time.DateOnly)
	job := &FetchNewsJob{
		Filters: types.NewFilteringParams("", d, "", ""),
	}

	err = job.Run()
	if err != nil {
		t.Fatalf("Run() returned an error: %v", err)
	}

	CwdPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	expectedFilename := filepath.Join(CwdPath, parsers.CmdDir,
		parsers.ParsersDir,
		parsers.DataDir,
		d)

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
}
