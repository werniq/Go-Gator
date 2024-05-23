package parsingInstructions

import (
	"newsAggr/cmd/types"
)

type InstructionsFactory interface {
	CreateApplyKeywordInstruction() ParsingInstruction
	CreateApplyDataRangeInstruction() ParsingInstruction
}

type ParsingInstruction interface {
	Apply(article types.News, params *types.ParsingParams) bool
}

type GoGatorInstructionFactory struct{}

func (g GoGatorInstructionFactory) CreateApplyKeywordInstruction() ParsingInstruction {
	return ApplyKeywordsInstruction{}
}

func (g GoGatorInstructionFactory) CreateApplyDataRangeInstruction() ParsingInstruction {
	return ApplyDateRangeInstruction{}
}
