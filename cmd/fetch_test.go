package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"log"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/templates"
	"newsaggr/cmd/types"
	"newsaggr/cmd/validator"
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

			checkFlagErr(tt.err)

			if got := buf.String(); got != tt.expectedLogs {
				t.Errorf("checkFlagErr() = %v, want %v", got, tt.expectedLogs)
			}
		})
	}
}

func TestAddFetchNewsCmd(t *testing.T) {
	fetchNews := AddFetchNewsCmd()

	runFunc := func(cmd *cobra.Command, args []string) {
		// retrieve optional parameters
		keywords, err := cmd.Flags().GetString(KeywordFlag)
		checkFlagErr(err)
		dateFrom, err := cmd.Flags().GetString(DateFromFlag)
		checkFlagErr(err)
		dateEnd, err := cmd.Flags().GetString(DateEndFlag)
		checkFlagErr(err)
		sources, err := cmd.Flags().GetString(SourcesFlag)
		checkFlagErr(err)

		// Validate user arguments
		if dateEnd != "" && dateFrom != "" {
			if dateFrom > dateEnd {
				log.Fatalln("Date from can not be after date end.")
			}
		}

		err = validator.ByDate(dateFrom)
		if err != nil {
			log.Fatalln("Error validating date: ", err)
		}

		err = validator.ByDate(dateEnd)
		if err != nil {
			log.Fatalln("Error validating date: ", err)
		}

		err = validator.BySources(sources)
		if err != nil {
			log.Fatalln("Error validating sources: ", err)
		}

		// Split and validate sources
		f := types.NewFilteringParams(keywords, dateFrom, dateEnd, sources)

		// parsing news by sources and applying params to those news
		news, err := parsers.ParseBySource(sources)
		if err != nil {
			log.Fatalln("Error parsing news: ", err)
		}

		news = filters.Apply(news, f)

		// output using go templates
		if err = templates.PrintTemplate(f, news); err != nil {
			log.Fatalln(err)
		}

		log.Println(len(news))
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
