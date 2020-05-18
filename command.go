package dgc

import "strings"

// Command represents a simple command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []string
	IgnoreCase  bool
	SubCommands []*Command
	Handler     CommandHandler
}

// CommandHandler represents a handler for a command
type CommandHandler func(*Ctx)

// trigger triggers the given command
func (command *Command) trigger(ctx *Ctx) {
	// Check if the first argument matches a sub command
	if len(ctx.Arguments.arguments) > 0 {
		argument := ctx.Arguments.Get(0)
		for _, subCommand := range command.SubCommands {
			valid := false
			if equals(argument.Raw(), subCommand.Name, subCommand.IgnoreCase) {
				valid = true
			} else {
				// Check if the first argument matches one of the aliases
				for _, alias := range subCommand.Aliases {
					if equals(argument.Raw(), alias, subCommand.IgnoreCase) {
						valid = true
					}
				}
			}

			if valid {
				// Define the arguments for the sub command
				arguments := ParseArguments("")
				if ctx.Arguments.Amount() > 1 {
					arguments = ParseArguments(strings.Join(strings.Split(ctx.Arguments.Raw(), " ")[1:], " "))
				}

				// Trigger the sub command
				subCommand.trigger(&Ctx{
					Session:       ctx.Session,
					Event:         ctx.Event,
					Arguments:     arguments,
					CustomObjects: ctx.CustomObjects,
					Router:        ctx.Router,
					Command:       subCommand,
				})
				return
			}
		}
	}

	// Run all middlewares assigned to this command
	for _, flag := range command.Flags {
		for _, middleware := range ctx.Router.Middlewares[flag] {
			if !middleware(ctx) {
				return
			}
		}
	}

	// Handle this command if the first argument matched no sub command
	command.Handler(ctx)
}
