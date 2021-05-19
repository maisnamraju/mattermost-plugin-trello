package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	botUserName    = "trello"
	botDisplayName = "Trello"
	botDescription = "Created by the Trello Plugin."

	autolinkPluginID = "mattermost-autolink"

	// Move these two to the plugin settings if admins need to adjust them.
	WebhookMaxProcsPerServer = 20
	WebhookBufferSize        = 10000
	PluginRepo               = "https://github.com/maisnamraju/mattermost-plugin-trello"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type TrelloPlugin struct {
	plugin.MattermostPlugin

	botID string

	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *Configuration

	router *mux.Router
}

func (p *TrelloPlugin) CreateBoard(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	var client = p.CreateClient(token)

	board := trello.NewBoard(r.URL.Query().Get("name"))
	board.Desc = r.URL.Query().Get("desc")

	// POST
	err := client.CreateBoard(&board, trello.Defaults())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, "Board Created")
}

func (p *TrelloPlugin) ListBoards(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	var client = p.CreateClient(token)

	boards, err := client.GetMyBoards()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(w, boards)
}

func (p *TrelloPlugin) CreateCard(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	var client = p.CreateClient(token)

	card := trello.Card{
		Name: r.URL.Query().Get("name"),
		Desc: r.URL.Query().Get("desc"),
	}
	err := client.CreateCard(&card, trello.Defaults())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(w, card)
}

func (p *TrelloPlugin) ListCards(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	var client = p.CreateClient(token)

	board, err := client.GetBoard(r.URL.Query().Get("id"), trello.Defaults())
	if err != nil {
		log.Fatalln(err)
	}
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		// Handle error
		log.Fatalln(err)
	}

	fmt.Println(w, cards)
}
