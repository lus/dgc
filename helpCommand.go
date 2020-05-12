package dgc

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// generalHelp handles the general help command
func generalHelp(ctx *Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		specificHelp(ctx)
		return
	}

	// Define useful variables
	channelID := ctx.Event.ChannelID
	session := ctx.Session

	// Send the general help embed
	embed, _ := renderDefaultGeneralHelpEmbed(ctx.Router, 1)
	message, _ := ctx.Session.ChannelMessageSendEmbed(channelID, embed)

	// Add the reactions to the message
	session.MessageReactionAdd(channelID, message.ID, "⬅️")
	session.MessageReactionAdd(channelID, message.ID, "❌")
	session.MessageReactionAdd(channelID, message.ID, "➡️")

	// Define the message as a help message
	ctx.Router.helpMessages[channelID+":"+message.ID] = 1
}

// specificHelp handles the specific help
func specificHelp(ctx *Ctx) {
	// Define the command name
	commandName := ctx.Arguments.Get(0).Raw()

	// Define the command
	var command *Command
	for _, cmd := range ctx.Router.Commands {
		if equals(commandName, cmd.Name, cmd.IgnoreCase) {
			command = cmd
			break
		}
	}

	// Send the help embed
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, renderDefaultSpecificHelpEmbed(ctx, command))
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func renderDefaultGeneralHelpEmbed(router *Router, page int) (*discordgo.MessageEmbed, int) {
	// Define useful variables
	commands := router.Commands
	prefix := router.Prefixes[0]

	// Calculate the amount of pages
	pageAmount := int(math.Ceil(float64(len(commands)) / 5))
	if page > pageAmount {
		page = pageAmount
	}
	if page <= 0 {
		page = 1
	}

	// Calculate the slice of commands to display on this page
	startingIndex := (page - 1) * 5
	endingIndex := startingIndex + 5
	if page == pageAmount {
		endingIndex = len(commands)
	}
	displayCommands := commands[startingIndex:endingIndex]

	// Prepare the fields for the embed
	fields := make([]*discordgo.MessageEmbedField, len(displayCommands))
	for index, command := range displayCommands {
		fields[index] = &discordgo.MessageEmbedField{
			Name:   command.Name,
			Value:  "`" + command.Description + "`",
			Inline: false,
		}
	}

	// Return the embed and the new page
	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command List (Page " + strconv.Itoa(page) + "/" + strconv.Itoa(pageAmount) + ")",
		Description: "These are all the available commands. Type `" + prefix + "help <command name>` to find out more about a specific command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields:      fields,
	}, page
}

// renderDefaultSpecificHelpEmbed renders the specific help embed of the given command
func renderDefaultSpecificHelpEmbed(ctx *Ctx, command *Command) *discordgo.MessageEmbed {
	// Define useful variables
	prefix := ctx.Router.Prefixes[0]

	// Check if the command is invalid
	if command == nil {
		return &discordgo.MessageEmbed{
			Type:      "rich",
			Title:     "Error",
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	// Define the aliases string
	aliases := "No aliases"
	if len(command.Aliases) > 0 {
		aliases = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	// Return the embed
	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command Information",
		Description: "Displaying the information for the `" + command.Name + "` command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
		},
	}
}
