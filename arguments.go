package dgc

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/karrick/tparse/v2"
)

var (
	// RegexArguments defines the regex the argument string has to match
	RegexArguments = regexp.MustCompile("(\"[^\"]+\"|[^\\s]+)")

	// RegexUserMention defines the regex a user mention has to match
	RegexUserMention = regexp.MustCompile("<@!?(\\d+)>")

	// RegexRoleMention defines the regex a role mention has to match
	RegexRoleMention = regexp.MustCompile("<@&(\\d+)>")

	// RegexChannelMention defines the regex a channel mention has to match
	RegexChannelMention = regexp.MustCompile("<#(\\d+)>")

	// RegexBigCodeblock defines the regex a big codeblock has to match
	RegexBigCodeblock = regexp.MustCompile("(?s)\\n*```(?:([\\w.\\-]*)\\n)?(.*)```")

	// RegexSmallCodeblock defines the regex a small codeblock has to match
	RegexSmallCodeblock = regexp.MustCompile("(?s)\\n*`(.*)`")

	// CodeblockLanguages defines which languages are valid codeblock languages
	CodeblockLanguages = []string{
		"as",
		"1c",
		"abnf",
		"accesslog",
		"actionscript",
		"ada",
		"ado",
		"adoc",
		"apache",
		"apacheconf",
		"applescript",
		"arduino",
		"arm",
		"armasm",
		"asciidoc",
		"aspectj",
		"atom",
		"autohotkey",
		"autoit",
		"avrasm",
		"awk",
		"axapta",
		"bash",
		"basic",
		"bat",
		"bf",
		"bind",
		"bnf",
		"brainfuck",
		"c",
		"c++",
		"cal",
		"capnp",
		"capnproto",
		"cc",
		"ceylon",
		"clean",
		"clj",
		"clojure-repl",
		"clojure",
		"cls",
		"cmake.in",
		"cmake",
		"cmd",
		"coffee",
		"coffeescript",
		"console",
		"coq",
		"cos",
		"cpp",
		"cr",
		"craftcms",
		"crm",
		"crmsh",
		"crystal",
		"cs",
		"csharp",
		"cson",
		"csp",
		"css",
		"d",
		"dart",
		"dcl",
		"delphi",
		"dfm",
		"diff",
		"django",
		"dns",
		"do",
		"docker",
		"dockerfile",
		"dos",
		"dpr",
		"dsconfig",
		"dst",
		"dts",
		"dust",
		"ebnf",
		"elixir",
		"elm",
		"erb",
		"erl",
		"erlang-repl",
		"erlang",
		"excel",
		"f90",
		"f95",
		"feature",
		"fix",
		"flix",
		"fortran",
		"freepascal",
		"fs",
		"fsharp",
		"gams",
		"gauss",
		"gcode",
		"gemspec",
		"gherkin",
		"glsl",
		"gms",
		"go",
		"golang",
		"golo",
		"gradle",
		"graph",
		"groovy",
		"gss",
		"gyp",
		"h",
		"h++",
		"haml",
		"handlebars",
		"haskell",
		"haxe",
		"hbs",
		"hpp",
		"hs",
		"hsp",
		"html.handlebars",
		"html.hbs",
		"html",
		"htmlbars",
		"http",
		"https",
		"hx",
		"hy",
		"hylang",
		"i7",
		"iced",
		"icl",
		"inform7",
		"ini",
		"instances",
		"irb",
		"irpf90",
		"java",
		"javascript",
		"jboss-cli",
		"jinja",
		"js",
		"json",
		"jsp",
		"jsx",
		"julia",
		"k",
		"kdb",
		"kotlin",
		"lasso",
		"lassoscript",
		"lazarus",
		"ldif",
		"leaf",
		"less",
		"lfm",
		"lisp",
		"livecodeserver",
		"livescript",
		"llvm",
		"lpr",
		"ls",
		"lsl",
		"lua",
		"m",
		"mak",
		"makefile",
		"markdown",
		"mathematica",
		"matlab",
		"maxima",
		"md",
		"mel",
		"mercury",
		"mips",
		"mipsasm",
		"mizar",
		"mk",
		"mkd",
		"mkdown",
		"ml",
		"mm",
		"mma",
		"mojolicious",
		"monkey",
		"moo",
		"moon",
		"moonscript",
		"n1ql",
		"nc",
		"nginx",
		"nginxconf",
		"nim",
		"nimrod",
		"nix",
		"nixos",
		"nsis",
		"obj-c",
		"objc",
		"objectivec",
		"ocaml",
		"openscad",
		"osascript",
		"oxygene",
		"p21",
		"parser3",
		"pas",
		"pascal",
		"patch",
		"pb",
		"pbi",
		"pcmk",
		"perl",
		"pf.conf",
		"pf",
		"php",
		"php3",
		"php4",
		"php5",
		"php6",
		"pl",
		"plist",
		"pm",
		"podspec",
		"pony",
		"powershell",
		"pp",
		"processing",
		"profile",
		"prolog",
		"protobuf",
		"ps",
		"puppet",
		"purebasic",
		"py",
		"python",
		"q",
		"qml",
		"qt",
		"r",
		"rb",
		"rib",
		"roboconf",
		"rs",
		"rsl",
		"rss",
		"ruby",
		"ruleslanguage",
		"rust",
		"scad",
		"scala",
		"scheme",
		"sci",
		"scilab",
		"scss",
		"sh",
		"shell",
		"smali",
		"smalltalk",
		"sml",
		"sqf",
		"sql",
		"st",
		"stan",
		"stata",
		"step",
		"step21",
		"stp",
		"styl",
		"stylus",
		"subunit",
		"sv",
		"svh",
		"swift",
		"taggerscript",
		"tao",
		"tap",
		"tcl",
		"tex",
		"thor",
		"thrift",
		"tk",
		"toml",
		"tp",
		"ts",
		"twig",
		"typescript",
		"v",
		"vala",
		"vb",
		"vbnet",
		"vbs",
		"vbscript-html",
		"vbscript",
		"verilog",
		"vhdl",
		"vim",
		"wildfly-cli",
		"x86asm",
		"xhtml",
		"xjb",
		"xl",
		"xls",
		"xlsx",
		"xml",
		"xpath",
		"xq",
		"xquery",
		"xsd",
		"xsl",
		"yaml",
		"yml",
		"zep",
		"zephir",
		"zone",
		"zsh",
	}
)

// Arguments represents the arguments that may be used in a command context
type Arguments struct {
	raw       string
	arguments []*Argument
}

// Codeblock represents a Discord codeblock
type Codeblock struct {
	Language string
	Content  string
}

// ParseArguments parses the raw string into several arguments
func ParseArguments(raw string) *Arguments {
	// Define the raw arguments
	rawArguments := RegexArguments.FindAllString(raw, -1)
	arguments := make([]*Argument, len(rawArguments))

	// Parse the raw arguments into correct ones
	for index, rawArgument := range rawArguments {
		rawArgument = stringTrimPreSuffix(rawArgument, "\"")
		arguments[index] = &Argument{
			raw: rawArgument,
		}
	}

	// Return the arguments structure
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
		return &Argument{
			raw: "",
		}
	}
	return arguments.arguments[n]
}

// Remove removes the n'th argument
func (arguments *Arguments) Remove(n int) {
	// Check if the given index is valid
	if arguments.Amount() <= n {
		return
	}

	// Set the new argument slice
	arguments.arguments = append(arguments.arguments[:n], arguments.arguments[n+1:]...)

	// Set the new raw string
	raw := ""
	for _, argument := range arguments.arguments {
		raw += argument.raw + " "
	}
	arguments.raw = strings.TrimSpace(raw)
}

// AsCodeblock parses the given arguments as a codeblock
func (arguments *Arguments) AsCodeblock() *Codeblock {
	raw := arguments.Raw()

	// Check if the raw string is a big codeblock
	matches := RegexBigCodeblock.MatchString(raw)
	if !matches {
		// Check if the raw string is a small codeblock
		matches = RegexSmallCodeblock.MatchString(raw)
		if matches {
			submatches := RegexSmallCodeblock.FindStringSubmatch(raw)
			return &Codeblock{
				Language: "",
				Content:  submatches[1],
			}
		}
		return nil
	}

	// Define the content and the language
	submatches := RegexBigCodeblock.FindStringSubmatch(raw)
	language := ""
	content := submatches[1] + submatches[2]
	if submatches[1] != "" && !stringArrayContains(CodeblockLanguages, submatches[1], false) {
		language = submatches[1]
		content = submatches[2]
	}

	// Return the codeblock
	return &Codeblock{
		Language: language,
		Content:  content,
	}
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

// AsDuration parses the given argument into a duration
func (argument *Argument) AsDuration() (time.Duration, error) {
	return tparse.AbsoluteDuration(time.Now(), argument.raw)
}
