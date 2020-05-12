package dgc

import "github.com/bwmarrin/discordgo"

// Router represents a DiscordGo command router
type Router struct {
	Prefixes    []string
	BotsAllowed bool
	commands    []*Command
}

// RegisterCmd registers a new command
func (router *Router) RegisterCmd(command *Command) {
	router.commands = append(router.commands, command)
}

// Initialize initializes the message event listener
func (router *Router) Initialize(session *discordgo.Session) {
	session.AddHandler(router.handler())
}

// handler provides the discordgo handler for the given router
func (router *Router) handler() func(*discordgo.Session, *discordgo.MessageCreate) {
	// TODO: Implement command parsing
}
