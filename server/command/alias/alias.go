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
	case "remove":
		return executeAliasRemove(args, fields, client, api)
	case "list":
		return executeAliasList(args, client)
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

func executeAliasRemove(args *model.CommandArgs, fields []string, client *pluginapi.Client, api plugin.API) (*model.CommandResponse, error) {
	if len(fields) < 3 || !strings.HasPrefix(fields[2], "@") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: /alias remove @username",
		}, nil
	}

	username := strings.TrimPrefix(fields[2], "@")
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

	if _, ok := aliases[targetUser.Id]; !ok {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Alias for @%s not found.", username),
		}, nil
	}

	delete(aliases, targetUser.Id)

	ok, appErr := client.KV.Set(storeKey, aliases)
	if !ok || appErr != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to update alias store.",
		}, nil
	}

	api.PublishWebSocketEvent("alias_update", nil, &model.WebsocketBroadcast{
		UserId: args.UserId,
	})

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Alias for @%s removed.", username),
	}, nil
}

func executeAliasList(args *model.CommandArgs, client *pluginapi.Client) (*model.CommandResponse, error) {
	storeKey := aliasKeyPrefix + args.UserId
	aliases := map[string]string{}
	_ = client.KV.Get(storeKey, &aliases)

	if len(aliases) == 0 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "You have no aliases set.",
		}, nil
	}

	var builder strings.Builder
	builder.WriteString("Your aliases:\n")
	for userID, alias := range aliases {
		user, err := client.User.Get(userID)
		if err != nil {
			continue
		}
		builder.WriteString(fmt.Sprintf("@%s → \"%s\"\n", user.Username, alias))
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         builder.String(),
	}, nil
}
