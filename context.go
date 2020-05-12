package dgc

import "github.com/bwmarrin/discordgo"

// Ctx represents the context for a command event
type Ctx struct {
	Session       *discordgo.Session
	Event         *discordgo.MessageCreate
	Arguments     *Arguments
	CustomObjects map[string]interface{}
}
