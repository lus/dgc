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
	PingHandler      CommandHandler
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

// RegisterDefaultHelpCommand registers the default help command
func (router *Router) RegisterDefaultHelpCommand(session *discordgo.Session) {
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
		IgnoreCase:  true,
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

		// Check if the message is empty after the prefix processing
		if content == "" {
			return
		}

		// Split the processed message
		split := strings.Split(content, " ")

		// Check if the message starts with a command name
		for _, command := range router.Commands {
			valid := false
			if equals(split[0], command.Name, command.IgnoreCase) {
				valid = true
			} else {
				// Check if the message starts with one of the aliases
				for _, alias := range command.Aliases {
					if equals(split[0], alias, command.IgnoreCase) {
						valid = true
					}
				}
			}

			if valid {
				// Define the arguments for the command
				arguments := ParseArguments("")
				if len(split) > 1 {
					arguments = ParseArguments(strings.Join(split[1:], " "))
				}

				// Define the context
				ctx := &Ctx{
					Session:       session,
					Event:         event,
					Arguments:     arguments,
					CustomObjects: map[string]interface{}{},
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
