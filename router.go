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
	Middlewares      []Middleware
	PingHandler      ExecutionHandler
	Storage          map[string]*ObjectsMap
}

// Create makes sure all maps get initialized
func Create(router *Router) *Router {
	router.Storage = make(map[string]*ObjectsMap)
	return router
}

// RegisterCmd registers a new command
func (router *Router) RegisterCmd(command *Command) {
	router.Commands = append(router.Commands, command)
}

// GetCmd returns the command with the given name if it exists
func (router *Router) GetCmd(name string) *Command {
	for _, command := range router.Commands {
		if equals(command.Name, name, command.IgnoreCase) || stringArrayContains(command.Aliases, name, command.IgnoreCase) {
			return command
		}
	}
	return nil
}

// RegisterMiddleware registers a new middleware
func (router *Router) RegisterMiddleware(middleware Middleware) {
	router.Middlewares = append(router.Middlewares, middleware)
}

// InitializeStorage initializes a storage map
func (router *Router) InitializeStorage(name string) {
	router.Storage[name] = newObjectsMap()
}

// Initialize initializes the message event listener
func (router *Router) Initialize(session *discordgo.Session) {
	session.AddHandler(router.Handler())
}

// Handler provides the discordgo handler for the given router
func (router *Router) Handler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, event *discordgo.MessageCreate) {
		// Define useful variables
		message := event.Message
		content := message.Content

		// Check if the message was sent by a bot
		if message.Author.Bot && !router.BotsAllowed {
			return
		}

		// Execute the ping handler if the message equals the current bot's mention
		if (content == "<@!"+session.State.User.ID+">" || content == "<@"+session.State.User.ID+">") && router.PingHandler != nil {
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
			var toCheck []string
			for _, alias := range command.Aliases {
				toCheck = append(toCheck, alias)
			}
			toCheck = append(toCheck, command.Name)

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

				// Trigger the command
				command.trigger(ctx)
			}
		}
	}
}
