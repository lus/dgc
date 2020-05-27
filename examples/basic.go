package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// Open a simple Discord session
	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	err = session.Open()
	if err != nil {
		panic(err)
	}

	// Wait for the user to cancel the process
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()

	// Create a dgc router
	// NOTE: The dgc.Create function makes sure all the maps get initialized
	router := dgc.Create(&dgc.Router{
		// We will allow '!' and 'example!' as the bot prefixes
		Prefixes: []string{
			"!",
			"example!",
		},

		// We will ignore the prefix case, so 'eXaMpLe!' is also a valid prefix
		IgnorePrefixCase: true,

		// We don't want bots to be able to execute our commands
		BotsAllowed: false,

		// We may initialize our commands in here, but we will use the corresponding method later on
		Commands: []*dgc.Command{},

		// We may inject our middlewares in here, but we will also use the corresponding method later on
		Middlewares: []dgc.Middleware{},

		// This handler gets called if the bot just got pinged (no argument provided)
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("Pong!")
		},
	})

	// Register the default help command
	router.RegisterDefaultHelpCommand(session, nil)

	// Register a simple middleware that injects a custom object
	router.RegisterMiddleware(func(next dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			// Inject a custom object into the context
			ctx.CustomObjects.Set("myObject", 69)

			// You can retrieve the object like this
			obj := ctx.CustomObjects.MustGet("myObject").(int)
			fmt.Println(obj)

			// Call the next execution handler
			next(ctx)
		}
	})

	// Register a simple command that responds with our custom object
	router.RegisterCmd(&dgc.Command{
		// We want to use 'obj' as the primary name of the command
		Name: "obj",

		// We also want the command to get triggered with the 'object' alias
		Aliases: []string{
			"object",
		},

		// These fields get displayed in the default help messages
		Description: "Responds with the injected custom object",
		Usage:       "obj",
		Example:     "obj",

		// You can assign custom flags to a command to use them in middlewares
		Flags: []string{},

		// We want to ignore the command case
		IgnoreCase: true,

		// You may define sub commands in here
		SubCommands: []*dgc.Command{},

		// We want the user to be able to execute this command once in five seconds and the cleanup interval shpuld be one second
		RateLimiter: dgc.NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *dgc.Ctx) {
			ctx.RespondText("You are being rate limited!")
		}),

		// Now we want to define the command handler
		Handler: objCommand,
	})
}

func objCommand(ctx *dgc.Ctx) {
	// Respond with the just set custom object
	ctx.RespondText(strconv.Itoa(ctx.CustomObjects.MustGet("myObject").(int)))
}
