package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUsername    = "trello"
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

	botUserID string

	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *Configuration
}

// OnActivate ensure the bot account exists
func (p *TrelloPlugin) OnActivate() error {
	bot := &model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: botDescription,
	}
	botUserID, appErr := p.Helpers.EnsureBot(bot)
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure bot user")
	}
	p.botUserID = botUserID

	err := p.API.RegisterCommand(getCommand())
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}
