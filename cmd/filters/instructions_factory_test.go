package filters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoGatorInstructionFactory_CreateApplyDataRangeInstruction(t *testing.T) {
	g := GoGatorInstructionFactory{}

	dataRangeInstruction := g.CreateApplyDataRangeInstruction()

	assert.Equal(t, dataRangeInstruction, ApplyDateRangeInstruction{})
}

func TestGoGatorInstructionFactory_CreateApplyKeywordInstruction(t *testing.T) {
	g := GoGatorInstructionFactory{}

	keywordInstruction := g.CreateApplyKeywordInstruction()

	assert.Equal(t, keywordInstruction, ApplyKeywordsInstruction{})
}
