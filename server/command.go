package main

import (
	"fmt"
	"strings"

	"github.com/adlio/trello"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const commandHelp = `use the following commands for performaning actions 
*|/trello create_board [name] [desc] | - create a board with description 
*|/trello list_boards | - list my boards
*|/trello create_card [name] [desc] | - create a card 
*|/trello list_cards [board_id] | - list cards in a board
`

const (
	trelloCreateBoardCommand = "create_board"
	trelloListBoards         = "list_boards"
	trelloCreateCard         = "create_card"
	trelloListCards          = "list_cards"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "trello",
		DisplayName:      "Trello",
		Description:      "Trello lets you create and list your boards and cards from trello",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: create_board,list_boards,create_card,list_cards",
		AutoCompleteHint: "[command]",
	}
}

func (p *TrelloPlugin) CreateBoard(boardName string, boardDesc string, args *model.CommandArgs) {
	var client = p.CreateClient()

	board := trello.NewBoard(boardName)
	board.Desc = boardDesc

	// POST
	err := client.CreateBoard(&board, trello.Defaults())
	if err != nil {
		p.postCommandResponse(args, "board name `%s` couldnt be created", boardName)
	}
	p.postCommandResponse(args, "board name `%s` created", boardName)
}

func (p *TrelloPlugin) ListBoards(args *model.CommandArgs) {
	var client = p.CreateClient()

	boards, err := client.GetMyBoards()
	if err != nil {
		p.postCommandResponse(args, "couldn't find boards")
	}
	if len(boards) == 0 {
		p.postCommandResponse(args, "couldn't find boards")
	}
	for _, board := range boards {
		p.postCommandResponse(args, "`%s` `%s`", board.Name, board.Desc)
	}
}

func (p *TrelloPlugin) CreateCard(cardName string, cardDesc string, args *model.CommandArgs) {
	var client = p.CreateClient()

	card := trello.Card{
		Name: cardName,
		Desc: cardDesc,
	}
	err := client.CreateCard(&card, trello.Defaults())
	if err != nil {
		p.postCommandResponse(args, "card name `%s` couldnt be created", cardName)
	}
	p.postCommandResponse(args, "card name `%s` created", cardName)
}

func (p *TrelloPlugin) ListCards(boardId string, args *model.CommandArgs) {
	var client = p.CreateClient()

	board, err := client.GetBoard(boardId)
	if err != nil {
		p.postCommandResponse(args, "cards couldnt be found for the board id")
	}
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		// Handle error
		p.postCommandResponse(args, "cards couldnt be found for the board id")
	}

	for _, card := range cards {
		p.postCommandResponse(args, "`%s` `%s`", card.Name, card.Desc)
	}
}

func (p *TrelloPlugin) postCommandResponse(args *model.CommandArgs, text string, textArgs ...interface{}) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf(text, textArgs...),
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

func (p *TrelloPlugin) validateCommand(action string, parameters []string) string {
	switch action {
	case trelloCreateBoardCommand:
		if len(parameters) < 2 {
			return "Please specify a name and description for the board"
		}
	case trelloListBoards:
		if len(parameters) > 0 {
			return "List command does not accept any extra parameters"
		}
	case trelloCreateCard:
		if len(parameters) < 2 {
			return "`Please specify a name and description for the card"
		}
	case trelloListCards:
		if len(parameters) == 1 {
			return "List cards needs the board id"
		}
	}
	return ""
}

func (p *TrelloPlugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	parameters := []string{}
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if command != "/trello" {
		return &model.CommandResponse{}, nil
	}

	if response := p.validateCommand(action, parameters); response != "" {
		p.postCommandResponse(args, response)
		return &model.CommandResponse{}, nil
	}

	switch action {
	case trelloCreateBoardCommand:
		boardName := parameters[0]
		boardDesc := parameters[1]
		p.CreateBoard(boardName, boardDesc, args)
		return &model.CommandResponse{}, nil
	case trelloListBoards:
		p.ListBoards(args)
		return &model.CommandResponse{}, nil
	case trelloCreateCard:
		cardName := parameters[0]
		cardDesc := parameters[1]
		p.CreateCard(cardName, cardDesc, args)
		return &model.CommandResponse{}, nil
	case trelloListCards:
		boardID := parameters[0]
		p.ListCards(boardID, args)
		return &model.CommandResponse{}, nil
	case "":
		text := "###### Mattermost Trello Commands Plugin - Slash Command Help\n" + strings.ReplaceAll(commandHelp, "|", "`")
		p.postCommandResponse(args, text)
		return &model.CommandResponse{}, nil
	}

	p.postCommandResponse(args, "Unknown action %v", action)
	return &model.CommandResponse{}, nil
}
