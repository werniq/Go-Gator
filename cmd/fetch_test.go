package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/templates"
	"gogator/cmd/types"
	"gogator/cmd/validator"
	"log"
	"reflect"
	"testing"
)

func TestCheckFlagErr(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedLogs string
	}{
		{
			name:         "Defined error message",
			err:          errors.New("flag accessed but not defined"),
			expectedLogs: "Unsupported flag: flag accessed but not defined\n",
		},
		{
			name:         "Undefined error message",
			err:          errors.New("some other error"),
			expectedLogs: "Error parsing flags: some other error\n",
		},
		{
			name:         "Nil error",
			err:          nil,
			expectedLogs: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil)

			validator.CheckFlagErr(tt.err)

			if got := buf.String(); got != tt.expectedLogs {
				t.Errorf("CheckFlagErr() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}

func TestAddFetchNewsCmd(t *testing.T) {
	fetchNews := FetchNewsCmd()

	runFunc := func(cmd *cobra.Command, args []string) {
		keywords, err := cmd.Flags().GetString(KeywordFlag)
		err = validator.CheckFlagErr(err)
		if err != nil {
			log.Fatalln(err)
		}

		dateFrom, err := cmd.Flags().GetString(DateFromFlag)
		err = validator.CheckFlagErr(err)
		if err != nil {
			log.Fatalln(err)
		}

		dateEnd, err := cmd.Flags().GetString(DateEndFlag)
		err = validator.CheckFlagErr(err)
		if err != nil {
			log.Fatalln(err)
		}

		sources, err := cmd.Flags().GetString(SourcesFlag)
		err = validator.CheckFlagErr(err)
		if err != nil {
			log.Fatalln(err)
		}

		if dateEnd != "" && dateFrom != "" {
			if dateFrom > dateEnd {
				log.Fatalln("Date from can not be after date end.")
			}
		}

		v := &validator.ArgValidator{}
		err = v.Validate(sources, dateFrom, dateEnd)
		if err != nil {
			log.Fatalln(err)
		}

		f := types.NewFilteringParams(keywords, dateFrom, dateEnd, sources)

		news, err := parsers.ParseBySource(sources)
		if err != nil {
			log.Fatalln("Error parsing news: ", err)
		}

		news = filters.Apply(news, f)

		err = templates.PrintTemplate(f, news)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Verify the command properties
	assert.Equal(t, "fetch", fetchNews.Use, "Command use should be 'fetch'")
	assert.Equal(t, "Fetching news from downloaded data", fetchNews.Short, "Command short description should match")
	assert.Contains(t, fetchNews.Long, "This command parses HTML, XML and JSON files sorts them by given arguments", "Command long description should contain 'This command parses HTML, XML and JSON files'")

	// Verify the flags
	assert.NotNil(t, fetchNews.Flags().Lookup("keywords"), "Flag 'keywords' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("date-from"), "Flag 'date-from' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("date-end"), "Flag 'date-end' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("sources"), "Flag 'sources' should be defined")
	reflect.DeepEqual(fetchNews.Run, runFunc)
}
