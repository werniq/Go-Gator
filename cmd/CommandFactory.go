package cmd

import "newsAggr/cmd/types"

type CommandFactory interface {
	CreateFetchNewsCommand() Command
}

type Command interface {
	Execute(params ParsingParams) []types.News
}

type GoGatorCommandFactory struct{}

func (g GoGatorCommandFactory) CreateFetchNewsCommand() {

}
