package filters

import (
	"newsAggr/cmd/types"
)

type FactoryInterface interface {
	CreateApplyKeywordInstruction() Instruction
	CreateApplyDataRangeInstruction() Instruction
}

// Instruction will be applied to article and return bool if it matches given instruction
type Instruction interface {
	Apply(article types.News, params *types.FilteringParams) bool
}

type InstructionFactory struct{}

// CreateApplyKeywordInstruction initializes keyword instruction.
// It is used to check if article contains given keywords
func (g InstructionFactory) CreateApplyKeywordInstruction() Instruction {
	return ApplyKeywordsInstruction{}
}

// CreateApplyDataRangeInstruction initializes date range instructions.
// It is used to check if article is published in given date range
func (g InstructionFactory) CreateApplyDataRangeInstruction() Instruction {
	return ApplyDateRangeInstruction{}
}
