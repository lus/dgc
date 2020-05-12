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
	PingHandler      CommandHandler
}

// RegisterCmd registers a new command
func (router *Router) RegisterCmd(command *Command) {
	router.Commands = append(router.Commands, command)
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

				// Trigger the command
				command.trigger(&Ctx{
					Session:   session,
					Event:     event,
					Arguments: arguments,
				})
			}
		}
	}
}
