package dgc

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// RegexUserMention defines the regex a user mention has to match
	RegexUserMention = regexp.MustCompile("<@!?(\\d+)>")

	// RegexRoleMention defines the regex a role mention has to match
	RegexRoleMention = regexp.MustCompile("<@&(\\d+)>")

	// RegexChannelMention defines the regex a channel mention has to match
	RegexChannelMention = regexp.MustCompile("<#(\\d+)>")
)

// Arguments represents the arguments that may be used in a command context
type Arguments struct {
	raw       string
	arguments []*Argument
}

// ParseArguments parses the raw string into several arguments
func ParseArguments(raw string) *Arguments {
	// Split the raw string and parse it into an array of arguments
	split := strings.Split(raw, " ")
	arguments := make([]*Argument, len(split))
	for key, value := range split {
		arguments[key] = &Argument{
			raw: value,
		}
	}

	// Return the arguments
	return &Arguments{
		raw:       raw,
		arguments: arguments,
	}
}

// Raw returns the raw string value of the arguments
func (arguments *Arguments) Raw() string {
	return arguments.raw
}

// AsSingle parses the given arguments as a single one
func (arguments *Arguments) AsSingle() *Argument {
	return &Argument{
		raw: arguments.raw,
	}
}

// Amount returns the amount of given arguments
func (arguments *Arguments) Amount() int {
	return len(arguments.arguments)
}

// Get returns the n'th argument
func (arguments *Arguments) Get(n int) *Argument {
	if arguments.Amount() <= n {
		return nil
	}
	return arguments.arguments[n]
}

// Argument represents a single argument
type Argument struct {
	raw string
}

// Raw returns the raw string value of the argument
func (argument *Argument) Raw() string {
	return argument.raw
}

// AsBool parses the given argument into a boolean
func (argument *Argument) AsBool() (bool, error) {
	return strconv.ParseBool(argument.raw)
}

// AsInt parses the given argument into an int32
func (argument *Argument) AsInt() (int, error) {
	return strconv.Atoi(argument.raw)
}

// AsInt64 parses the given argument into an int64
func (argument *Argument) AsInt64() (int64, error) {
	return strconv.ParseInt(argument.raw, 10, 64)
}

// AsUserMentionID returns the ID of the mentioned user or an empty string if it is no mention
func (argument *Argument) AsUserMentionID() string {
	// Check if the argument is a user mention
	matches := RegexUserMention.MatchString(argument.raw)
	if !matches {
		return ""
	}

	// Parse the user ID
	userID := RegexUserMention.FindStringSubmatch(argument.raw)[1]
	return userID
}

// AsRoleMentionID returns the ID of the mentioned role or an empty string if it is no mention
func (argument *Argument) AsRoleMentionID() string {
	// Check if the argument is a role mention
	matches := RegexRoleMention.MatchString(argument.raw)
	if !matches {
		return ""
	}

	// Parse the role ID
	roleID := RegexRoleMention.FindStringSubmatch(argument.raw)[1]
	return roleID
}

// AsChannelMentionID returns the ID of the mentioned channel or an empty string if it is no mention
func (argument *Argument) AsChannelMentionID() string {
	// Check if the argument is a channel mention
	matches := RegexChannelMention.MatchString(argument.raw)
	if !matches {
		return ""
	}

	// Parse the channel ID
	channelID := RegexChannelMention.FindStringSubmatch(argument.raw)[1]
	return channelID
}
