package cmd

import (
	"newsAggr/cmd/types"
)

type ParsingInstructionsFactory interface {
	CreateApplyKeywordInstruction() ParsingInstruction
	CreateApplyDataRangeInstruction() ParsingInstruction
}

type ParsingInstruction interface {
	Apply(article types.News, params ParsingParams) bool
}

type GoGatorInstructionFactory struct{}

func (g GoGatorInstructionFactory) CreateApplyKeywordInstruction() ParsingInstruction {
	return ApplyKeywordsInstruction{}
}

func (g GoGatorInstructionFactory) CreateApplyDataRangeInstruction() ParsingInstruction {
	return ApplyDateRangeInstruction{}
}
