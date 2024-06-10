package filters

import (
	"newsAggr/cmd/types"
)

type InstructionsFactory interface {
	CreateApplyKeywordInstruction() ParsingInstruction
	CreateApplyDataRangeInstruction() ParsingInstruction
}

// ParsingInstruction will be applied to article and return bool if it matches given instruction
type ParsingInstruction interface {
	Apply(article types.News, params *types.FilteringParams) bool
}

type GoGatorInstructionFactory struct{}

// CreateApplyKeywordInstruction initializes keyword instruction.
// It is used to check if article contains given keywords
func (g GoGatorInstructionFactory) CreateApplyKeywordInstruction() ParsingInstruction {
	return ApplyKeywordsInstruction{}
}

// CreateApplyDataRangeInstruction initializes date range instructions.
// It is used to check if article is published in given date range
func (g GoGatorInstructionFactory) CreateApplyDataRangeInstruction() ParsingInstruction {
	return ApplyDateRangeInstruction{}
}
