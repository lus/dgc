[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/Lukaesebrot/dgc)

# dgc

A DiscordGo command router with tons of useful features
If you find any bugs or if you have a feature request, please tell me using an issue.

## Deprecation notice

After thinking a lot about how to continue this project, I decided I won't.
The main reason for this is that Discord will make the message content intent privileged in 2022 and thus will force every bot developer to use slash commands.
I don't feel motivated to continue a project which will become pretty much redundant in the next months, especially because a major rewrite of some features would be mandatory for me as I started the project at the beginning of my Go learning period and I evolved a lot since then.
I think I could use the time much better for new, cooler projects.

## Basic example

This just shows a very basic example:
```go
func main() {
    // Discord bot logic here
    session, _ := ...

    // Create a new command router
    router := dgc.Create(&dgc.Router{
        Prefixes: []string{"!"},
    })

    // Register a simple ping command
    router.RegisterCmd(&dgc.Command{
        Name: "ping",
        Description: "Responds with 'pong!'",
        Usage: "ping",
        Example: "ping",
        IgnoreCase: true,
        Handler: func(ctx *dgc.Ctx) {
            ctx.RespondText("Pong!")
        },
    })

    // Initialize the router
    router.Initialize(session)
}
```

## Usage

You can find examples for the more complex usage and all the integrated features in the `examples/*.go` files.

## Arguments

To find out how you can use the integrated argument parser, just look into the `examples/arguments.go` file.
