[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/Lukaesebrot/dgc)

# dgc
A DiscordGo command router with tons of useful features
If you find any bugs or if you have a feature request, please tell me using an issue.


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
