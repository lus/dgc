package dgc

import "github.com/bwmarrin/discordgo"

// Ctx represents the context for a command event
type Ctx struct {
	Session       *discordgo.Session
	Event         *discordgo.MessageCreate
	Arguments     *Arguments
	CustomObjects *ObjectsMap
	Router        *Router
	Command       *Command
}

// ExecutionHandler represents a handler for a context execution
type ExecutionHandler func(*Ctx)

// RespondText responds with the given text message
func (ctx *Ctx) RespondText(text string) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, text)
	return err
}

// RespondEmbed responds with the given embed message
func (ctx *Ctx) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embed)
	return err
}

// RespondTextEmbed responds with the given text and embed message
func (ctx *Ctx) RespondTextEmbed(text string, embed *discordgo.MessageEmbed) error {
	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
		Content: text,
		Embed:   embed,
	})
	return err
}
