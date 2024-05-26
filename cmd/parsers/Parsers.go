package parsers

import "newsAggr/cmd/types"

type Parsers interface {
	Parse(params *types.FilteringParams) []types.News
}
