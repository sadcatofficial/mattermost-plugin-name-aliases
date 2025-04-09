package alias

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

const aliasKeyPrefix = "alias_store_"

func ExecuteAliasCommand(args *model.CommandArgs, client *pluginapi.Client, api plugin.API) (*model.CommandResponse, error) {
	fields := strings.Fields(args.Command)

	if len(fields) < 2 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: /alias set @username \"Alias\", /alias remove @username, /alias list",
		}, nil
	}

	switch fields[1] {
	case "set":
		return executeAliasSet(args, fields, client, api)
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown subcommand: %s", fields[1]),
		}, nil
	}
}

func executeAliasSet(args *model.CommandArgs, fields []string, client *pluginapi.Client, api plugin.API) (*model.CommandResponse, error) {
	if len(fields) < 4 || !strings.HasPrefix(fields[2], "@") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: /alias set @username \"Alias\"",
		}, nil
	}

	username := strings.TrimPrefix(fields[2], "@")
	alias := strings.Join(fields[3:], " ")
	alias = strings.Trim(alias, "\"")

	targetUser, err := client.User.GetByUsername(username)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("User @%s not found", username),
		}, nil
	}

	storeKey := aliasKeyPrefix + args.UserId

	aliases := map[string]string{}
	_ = client.KV.Get(storeKey, &aliases)

	aliases[targetUser.Id] = alias

	ok, appErr := client.KV.Set(storeKey, aliases)
	if !ok || appErr != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to store alias.",
		}, nil
	}

	api.PublishWebSocketEvent("alias_update", nil, &model.WebsocketBroadcast{
		UserId: args.UserId,
	})

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Alias for @%s set to \"%s\".", username, alias),
	}, nil
}
