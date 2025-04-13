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
		AutoCompleteDesc: "Manage display name aliases for users",
		AutoCompleteHint: "[set|remove|list]",
		AutocompleteData: buildAliasAutocompleteData(),
	})
	if err != nil {
		client.Log.Error("Failed to register command", "error", err)
	}
	return &Handler{
		client: client,
		api:    api,
	}
}

func buildAliasAutocompleteData() *model.AutocompleteData {
	root := model.NewAutocompleteData(aliasCommandTrigger, "[set|remove|list]", "Manage your display name aliases")

	set := model.NewAutocompleteData("set", "@username \"Alias\"", "Set a new alias for a user")
	remove := model.NewAutocompleteData("remove", "@username", "Remove an alias for a user")
	list := model.NewAutocompleteData("list", "", "List all your aliases")

	root.AddCommand(set)
	root.AddCommand(remove)
	root.AddCommand(list)

	return root
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
