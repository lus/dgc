package dgc

import (
	"strings"
)

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
	RateLimiter RateLimiter
	Handler     ExecutionHandler
}

// GetSubCmd returns the sub command with the given name if it exists
func (command *Command) GetSubCmd(name string) *Command {
	// Loop through all commands to find the correct one
	for _, subCommand := range command.SubCommands {
		// Define the slice to check
		toCheck := make([]string, len(subCommand.Aliases)+1)
		toCheck = append(toCheck, subCommand.Name)
		toCheck = append(toCheck, subCommand.Aliases...)

		// Check the prefix of the string
		if stringArrayContains(toCheck, name, subCommand.IgnoreCase) {
			return subCommand
		}
	}
	return nil
}

// NotifyRateLimiter notifies the rate limiter about a new execution and returns false if the user is being rate limited
func (command *Command) NotifyRateLimiter(ctx *Ctx) bool {
	if command.RateLimiter == nil {
		return true
	}
	return command.RateLimiter.NotifyExecution(ctx)
}

// trigger triggers the given command
func (command *Command) trigger(ctx *Ctx) {
	// Check if the first argument matches a sub command
	if len(ctx.Arguments.arguments) > 0 {
		argument := ctx.Arguments.Get(0).Raw()
		subCommand := command.GetSubCmd(argument)
		if subCommand != nil {
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

	// Prepare all middlewares
	nextHandler := command.Handler
	for _, middleware := range ctx.Router.Middlewares {
		nextHandler = middleware(nextHandler)
	}

	// Run all middlewares
	nextHandler(ctx)
}
