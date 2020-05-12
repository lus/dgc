package dgc

import (
	"strings"
)

// Command represents a simple command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	IgnoreCase  bool
	SubCommands []*Command
	Handler     CommandHandler
}

// CommandHandler represents a handler for a command
type CommandHandler func(ctx *Ctx)

// trigger triggers the given command
func (command *Command) trigger(ctx *Ctx) {
	// Check if the first argument matches a sub command
	for _, argument := range ctx.Arguments.arguments {
		for _, subCommand := range command.SubCommands {
			if equals(argument.Raw(), subCommand.Name, subCommand.IgnoreCase) {
				// Define the arguments for the sub command
				arguments := ParseArguments("")
				if ctx.Arguments.Amount() > 1 {
					arguments = ParseArguments(strings.Join(strings.Split(ctx.Arguments.Raw(), " ")[1:], " "))
				}

				// Trigger the sub command
				subCommand.trigger(&Ctx{
					Session:   ctx.Session,
					Event:     ctx.Event,
					Arguments: arguments,
				})
				return
			}
		}
	}

	// Handle this command if the first argument matched no sub command
	command.Handler(ctx)
}
