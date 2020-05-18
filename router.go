package dgc

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Router represents a DiscordGo command router
type Router struct {
	Prefixes         []string
	IgnorePrefixCase bool
	BotsAllowed      bool
	Commands         []*Command
	Middlewares      map[string][]Middleware
	PingHandler      ExecutionHandler
	helpMessages     map[string]int
}

// Middleware defines how a middleware looks like
type Middleware func(*Ctx) bool

// Create creates a new router and makes sure that all maps get initialized
func Create(router *Router) *Router {
	if router.Middlewares == nil {
		router.Middlewares = make(map[string][]Middleware)
	}
	router.helpMessages = make(map[string]int)
	return router
}

// RegisterCmd registers a new command
func (router *Router) RegisterCmd(command *Command) {
	router.Commands = append(router.Commands, command)
}

// GetCmd returns the command with the given name if it exists
func (router *Router) GetCmd(name string) *Command {
	for _, command := range router.Commands {
		if command.Name == name || stringArrayContains(command.Aliases, name, command.IgnoreCase) {
			return command
		}
	}
	return nil
}

// RegisterDefaultHelpCommand registers the default help command
func (router *Router) RegisterDefaultHelpCommand(session *discordgo.Session, rateLimiter *RateLimiter) {
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
		page := router.helpMessages[channelID+":"+messageID]
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
		router.helpMessages[channelID+":"+messageID] = page
	})

	// Register the default help command
	router.RegisterCmd(&Command{
		Name:        "help",
		Description: "Lists all the available commands or displays some information about a specific command",
		Usage:       "help [command name]",
		Example:     "help yourCommand",
		IgnoreCase:  true,
		RateLimiter: rateLimiter,
		Handler:     generalHelp,
	})
}

// AddMiddleware adds a new middleware
func (router *Router) AddMiddleware(flag string, middleware Middleware) {
	router.Middlewares[flag] = append(router.Middlewares[flag], middleware)
}

// Initialize initializes the message event listener
func (router *Router) Initialize(session *discordgo.Session) {
	session.AddHandler(router.handler())
}

// handler provides the discordgo handler for the given router
func (router *Router) handler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, event *discordgo.MessageCreate) {
		// Define useful variables
		message := event.Message
		content := message.Content

		// Check if the message was sent by a bot
		if message.Author.Bot && !router.BotsAllowed {
			return
		}

		// Execute the ping handler if the message equals the current bot's mention
		if content == "<@!"+session.State.User.ID+">" && router.PingHandler != nil {
			router.PingHandler(&Ctx{
				Session:   session,
				Event:     event,
				Arguments: ParseArguments(""),
				Router:    router,
			})
			return
		}

		// Check if the message starts with one of the defined prefixes
		hasPrefix, content := stringHasPrefix(content, router.Prefixes, router.IgnorePrefixCase)
		if !hasPrefix {
			return
		}

		// Get rid of additional spaces
		content = strings.Trim(content, " ")

		// Check if the message is empty after the prefix processing
		if content == "" {
			return
		}

		// Check if the message starts with a command name
		for _, command := range router.Commands {
			// Define an array containing the commands name and the aliases
			toCheck := make([]string, len(command.Aliases)+1)
			toCheck[0] = command.Name
			for index, alias := range command.Aliases {
				toCheck[index+1] = alias
			}

			// Check if the content is the current command
			isCommand, content := stringHasPrefix(content, toCheck, command.IgnoreCase)
			if !isCommand {
				continue
			}

			// Check if the remaining string is empty or starts with a space or newline
			isValid, content := stringHasPrefix(content, []string{" ", "\n"}, false)
			if content == "" || isValid {
				// Define the command context
				ctx := &Ctx{
					Session:       session,
					Event:         event,
					Arguments:     ParseArguments(content),
					CustomObjects: newObjectsMap(),
					Router:        router,
					Command:       command,
				}

				// Run all wildcard middlewares
				for _, middleware := range router.Middlewares["*"] {
					if !middleware(ctx) {
						return
					}
				}

				// Trigger the command
				command.trigger(ctx)
			}
		}
	}
}
