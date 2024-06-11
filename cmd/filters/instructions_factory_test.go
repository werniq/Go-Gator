package filters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoGatorInstructionFactory_CreateApplyDataRangeInstruction(t *testing.T) {
	g := InstructionFactory{}

	dataRangeInstruction := g.CreateApplyDataRangeInstruction()

	assert.Equal(t, dataRangeInstruction, ApplyDateRangeInstruction{})
}

func TestGoGatorInstructionFactory_CreateApplyKeywordInstruction(t *testing.T) {
	g := InstructionFactory{}

	keywordInstruction := g.CreateApplyKeywordInstruction()

	assert.Equal(t, keywordInstruction, ApplyKeywordsInstruction{})
}
