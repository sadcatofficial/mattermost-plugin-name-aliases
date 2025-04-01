package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/sadcatofficial/mattermost-plugin-name-aliases/server/command/alias"
)

type Handler struct {
	client *pluginapi.Client
}

type Command interface {
	Handle(args *model.CommandArgs) (*model.CommandResponse, error)
	executeHelloCommand(args *model.CommandArgs) *model.CommandResponse
}

const aliasCommandTrigger = "alias"

func (c *Handler) executeAliasCommand(args *model.CommandArgs) (*model.CommandResponse, error) {
	fields := strings.Fields(args.Command)

	if len(fields) < 4 || fields[1] != "set" || !strings.HasPrefix(fields[2], "@") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: /alias set @username \"Display Name\"",
		}, nil
	}

	username := strings.TrimPrefix(fields[2], "@")
	alias := strings.Join(fields[3:], " ")
	alias = strings.Trim(alias, "\"")

	user, err := c.client.User.GetByUsername(username)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("User @%s not found", username),
		}, nil
	}

	patch := &model.UserPatch{
		Nickname: &alias,
	}

	_, err = c.client.User.Patch(user.Id, patch)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to set alias.",
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Alias for @%s set to \"%s\".", username, alias),
	}, nil
}


// Register all your slash commands in the NewCommandHandler function.
func NewCommandHandler(client *pluginapi.Client) Command {
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
	}
}

// ExecuteCommand hook calls this method to execute the commands that were registered in the NewCommandHandler function.
func (c *Handler) Handle(args *model.CommandArgs) (*model.CommandResponse, error) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")

	switch trigger {
	case aliasCommandTrigger:
		return alias.ExecuteAliasCommand(args, c.client)
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s", args.Command),
		}, nil
	}
}

func (c *Handler) executeHelloCommand(args *model.CommandArgs) *model.CommandResponse {
	if len(strings.Fields(args.Command)) < 2 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Please specify a username",
		}
	}
	username := strings.Fields(args.Command)[1]
	return &model.CommandResponse{
		Text: "Hello, " + username,
	}
}
