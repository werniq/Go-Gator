package cmd

import (
	"newsAggr/cmd/types"
)

type ParsingInstructionsFactory interface {
	//ApplyKeywordInstruction(article types.News, pattern string) bool
	//ApplyStartingTimestamp(article types.News, pattern string) bool
	//ApplyEndingTimestamp(article types.News, pattern string) bool
	//ApplySource(article types.News, pattern string) bool

	CreateApplyKeywordInstruction() ParsingInstruction
	CreateApplyStartingTimestampInstruction() ParsingInstruction
	CreateApplyEndingTimestampInstruction() ParsingInstruction
	CreateApplySourceInstruction() ParsingInstruction
}

type ParsingInstruction interface {
	Apply(article types.News, pattern string) bool
}

type GoGatorInstructionFactory struct{}

func (g GoGatorInstructionFactory) CreateApplyStartingTimestampInstruction() ParsingInstruction {
	return ApplyStartingTimestampInstruction{}
}

func (g GoGatorInstructionFactory) CreateApplyEndingTimestampInstruction() ParsingInstruction {
	return ApplyStartingTimestampInstruction{}
}

func (g GoGatorInstructionFactory) CreateApplyKeywordInstruction() ParsingInstruction {
	return ApplyKeywordsInstruction{}
}

func (g GoGatorInstructionFactory) CreateApplySourceInstruction() ParsingInstruction {
	return ApplySourceInstruction{}
}
