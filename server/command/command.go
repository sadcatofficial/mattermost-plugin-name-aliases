package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/sadcatofficial/mattermost-plugin-name-aliases/server/command/alias"
)

type Handler struct {
	client *pluginapi.Client
	api    plugin.API
}

type Command interface {
	Handle(args *model.CommandArgs) (*model.CommandResponse, error)
}

const aliasCommandTrigger = "alias"

// Register all your slash commands in the NewCommandHandler function.
func NewCommandHandler(client *pluginapi.Client, api plugin.API) Command {
	err := client.SlashCommand.Register(&model.Command{
		Trigger:          aliasCommandTrigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Set a display name alias for a user",
		AutoCompleteHint: "set @username \"Alias\"",
		AutocompleteData: model.NewAutocompleteData(aliasCommandTrigger, "set @username \"Alias\"", "Set alias"),
	})
	if err != nil {
		client.Log.Error("Failed to register command", "error", err)
	}
	return &Handler{
		client: client,
		api:    api,
	}
}

// ExecuteCommand hook calls this method to execute the commands that were registered in the NewCommandHandler function.
func (c *Handler) Handle(args *model.CommandArgs) (*model.CommandResponse, error) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")

	switch trigger {
	case aliasCommandTrigger:
		return alias.ExecuteAliasCommand(args, c.client, c.api)
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s", args.Command),
		}, nil
	}
}
