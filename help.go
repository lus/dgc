package dgc

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// RegisterDefaultHelpCommand registers the default help command
func (router *Router) RegisterDefaultHelpCommand(session *discordgo.Session, rateLimiter RateLimiter) {
	// Initialize the helo messages storage
	router.InitializeStorage("dgc_helpMessages")

	// Initialize the reaction add listener
	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		// Define useful variables
		channelID := event.ChannelID
		messageID := event.MessageID
		userID := event.UserID

		// Check whether or not the reaction was added by the bot itself
		if event.UserID == session.State.User.ID {
			return
		}

		// Check whether or not the message is a help message
		rawPage, ok := router.Storage["dgc_helpMessages"].Get(channelID + ":" + messageID + ":" + event.UserID)
		if !ok {
			return
		}
		page := rawPage.(int)
		if page <= 0 {
			return
		}

		// Check which reaction was added
		reactionName := event.Emoji.Name
		switch reactionName {
		case "⬅️":
			// Update the help message
			embed, newPage := renderDefaultGeneralHelpEmbed(router, page-1)
			page = newPage
			session.ChannelMessageEditEmbed(channelID, messageID, embed)

			// Remove the reaction
			session.MessageReactionRemove(channelID, messageID, reactionName, userID)
			break
		case "❌":
			// Delete the help message
			session.ChannelMessageDelete(channelID, messageID)
			break
		case "➡️":
			// Update the help message
			embed, newPage := renderDefaultGeneralHelpEmbed(router, page+1)
			page = newPage
			session.ChannelMessageEditEmbed(channelID, messageID, embed)

			// Remove the reaction
			session.MessageReactionRemove(channelID, messageID, reactionName, userID)
			break
		}

		// Update the stores page
		router.Storage["dgc_helpMessages"].Set(channelID+":"+messageID+":"+event.UserID, page)
	})

	// Register the default help command
	router.RegisterCmd(&Command{
		Name:        "help",
		Description: "Lists all the available commands or displays some information about a specific command",
		Usage:       "help [command name]",
		Example:     "help yourCommand",
		IgnoreCase:  true,
		RateLimiter: rateLimiter,
		Handler:     generalHelpCommand,
	})
}

// generalHelpCommand handles the general help command
func generalHelpCommand(ctx *Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		specificHelpCommand(ctx)
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
	ctx.Router.Storage["dgc_helpMessages"].Set(channelID+":"+message.ID+":"+ctx.Event.Author.ID, 1)
}

// specificHelpCommand handles the specific help command
func specificHelpCommand(ctx *Ctx) {
	// Define the command names
	commandNames := strings.Split(ctx.Arguments.Raw(), " ")

	// Define the command
	var command *Command
	for index, commandName := range commandNames {
		if index == 0 {
			command = ctx.Router.GetCmd(commandName)
			continue
		}
		if command == nil {
			break
		}
		command = command.GetSubCmd(commandName)
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
				{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	// Define the sub commands string
	subCommands := "No sub commands"
	if len(command.SubCommands) > 0 {
		subCommandNames := make([]string, len(command.SubCommands))
		for index, subCommand := range command.SubCommands {
			subCommandNames[index] = subCommand.Name
		}
		subCommands = "`" + strings.Join(subCommandNames, "`, `") + "`"
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
			{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			{
				Name:   "Sub Commands",
				Value:  subCommands,
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
			{
				Name:   "Example",
				Value:  "```" + prefix + command.Example + "```",
				Inline: false,
			},
		},
	}
}
