package dgc

// Command represents a simple command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	SubCommands []*Command
}
