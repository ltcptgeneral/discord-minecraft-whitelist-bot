package app

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "whitelist-add",
			Description: "Add yourself to the whitelist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "minecraft-username",
					Description: "Minecraft Username",
					Required:    true,
				},
			},
		},
		{
			Name:        "whitelist-remove",
			Description: "Remove yourself from the whitelist",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
		{
			Name:        "whitelist-show",
			Description: "Display your whitelist link",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whitelist-add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			requestedUsername := (optionMap["minecraft-username"]).StringValue()
			invokingUserName := i.Member.User.Username
			invokingUserID := i.Member.User.ID

			// check requestedUsername validity
			valid := checkMcUsernameValid(requestedUsername)
			if !valid {
				simpleResponse(s, i, fmt.Sprintf("%s is not a valid minecraft username", requestedUsername))
				return
			}

			existingMCUsername, ok := db[invokingUserID]
			if !ok { // invoking user does not currently have an mc username associated
				// execute RCON
				response, err := conn.Execute(fmt.Sprintf("whitelist add %s", requestedUsername))
				if err != nil {
					simpleResponse(s, i, fmt.Sprintf("Failed to add %s to whitelist: %s", requestedUsername, err.Error()))
					return
				}

				// this can happen if the username is not a real player, OR if the user is already on the whitelist
				// In both cases we want to exit early to avoid adding an invalid username to the db
				if !strings.EqualFold(response, fmt.Sprintf("Added %s to the whitelist", requestedUsername)) {
					simpleResponse(s, i, fmt.Sprintf("Failed to add %s to whitelist: %s", requestedUsername, response))
					return
				}

				// save state to db
				db[invokingUserID] = requestedUsername
				SaveDB(dbPath, db)

				// send response
				simpleResponse(s, i, fmt.Sprintf("%s linked minecraft username %s", invokingUserName, requestedUsername))
			} else {
				simpleResponse(s, i, fmt.Sprintf("%s already has a linked minecraft username %s", invokingUserName, existingMCUsername))
			}

		},
		"whitelist-remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			invokingUserName := i.Member.User.Username
			invokingUserID := i.Member.User.ID

			existingMCUsername, ok := db[invokingUserID]
			if !ok { // invoking user does not currently have an mc username associated
				simpleResponse(s, i, fmt.Sprintf("%s does not have a linked minecraft username", invokingUserName))
			} else {
				// execute RCON
				response, err := conn.Execute(fmt.Sprintf("whitelist remove %s", existingMCUsername))
				if err != nil {
					simpleResponse(s, i, fmt.Sprintf("Failed to remove %s from whitelist: %s", existingMCUsername, err.Error()))
					return
				}

				// this can happen if the username is not a real player, OR if the user is already not on the whitelist
				if !strings.EqualFold(response, fmt.Sprintf("Removed %s from the whitelist", existingMCUsername)) {
					simpleResponse(s, i, fmt.Sprintf("Failed to remove %s from whitelist: %s", existingMCUsername, response))
					return
				}

				// save state to db
				delete(db, invokingUserID)
				SaveDB(dbPath, db)

				// send response
				simpleResponse(s, i, fmt.Sprintf("%s unlinked minecraft username %s", invokingUserName, existingMCUsername))

			}
		},
		"whitelist-show": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			invokingUserName := i.Member.User.Username
			invokingUserID := i.Member.User.ID

			existingMCUsername, ok := db[invokingUserID]

			if ok {
				simpleResponse(s, i, fmt.Sprintf("%s is %s", invokingUserName, existingMCUsername))
			} else {
				simpleResponse(s, i, fmt.Sprintf("%s does not have a linked minecraft username", invokingUserName))
			}
		},
	}
)
