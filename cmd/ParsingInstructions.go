package cmd

import (
	"log"
	"newsAggr/cmd/types"
	"strings"
	"time"
)

// containsKeywordsInText returns true if we have keyword in text
func containsKeywordsInText(text, keywordsSequence string) bool {
	keywords := strings.Split(keywordsSequence, " ")

	for _, keyword := range keywords {
		if containsPatternInText(text, keyword) ||
			containsPatternInText(text, keyword) {
			return true
		}
	}

	return false
}

// containsPatternInText checks if existing pattern is in the string
func containsPatternInText(text, pattern string) bool {
	if strings.Contains(text, pattern) {
		return true
	}

	return false
}

type ApplyStartingTimestampInstruction struct{}

func (a ApplyStartingTimestampInstruction) Apply(article types.News, timestamp string) bool {
	if article.PubDate == "" || timestamp == "" {
		return false
	}
	var publishedDate time.Time
	var startingTimestamp time.Time
	var err error

	publishedDate, err = time.Parse("2006-01-02", article.PubDate[:10])
	if err != nil {
		log.Fatalf("Error parsing timestamp: %v\n", err)
	}

	startingTimestamp, err = time.Parse("2006-01-02", timestamp[:10])
	if err != nil {
		log.Fatalf("Error parsing timestamp: %v\n", err)
	}

	if publishedDate.Before(startingTimestamp) {
		return true
	}

	return false
}

type ApplyEndingTimestampInstruction struct{}

func (a ApplyEndingTimestampInstruction) Apply(article types.News, timestamp string) bool {
	if article.PubDate == "" || timestamp == "" {
		return false
	}
	var publishedDate time.Time
	var startingTimestamp time.Time
	var err error

	publishedDate, err = time.Parse("2006-01-02", article.PubDate[:10])
	if err != nil {
		log.Fatalf("Error parsing timestamp: %v\n", err)
	}

	startingTimestamp, err = time.Parse("2006-01-02", timestamp[:10])
	if err != nil {
		log.Fatalf("Error parsing timestamp: %v\n", err)
	}

	if publishedDate.Before(startingTimestamp) {
		return false
	}

	return true
}

type ApplySourceInstruction struct{}

func (a ApplySourceInstruction) Apply(article types.News, timestamp string) bool {
	return true
}

type ApplyKeywordsInstruction struct{}

func (a ApplyKeywordsInstruction) Apply(article types.News, keywordsSequence string) bool {
	return containsKeywordsInText(article.Title, keywordsSequence) ||
		containsKeywordsInText(article.Description, keywordsSequence)
}
