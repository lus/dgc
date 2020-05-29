package main

import (
	"fmt"

	"github.com/Lukaesebrot/dgc"
)

// This example shows how to use the integrated argument parser

func someCommandHandler(ctx *dgc.Ctx) {
	// First of all, get the arguments object
	arguments := ctx.Arguments

	// Print the amount of arguments into the console
	amount := arguments.Amount()
	fmt.Println("Amount:", amount)

	// Print the raw argument string into the console
	raw := arguments.Raw()
	fmt.Println("Raw:", raw)

	// Parse it into a codeblock struct
	codeblock := arguments.AsCodeblock()
	if codeblock == nil {
		// Arguments aren't a codeblock
	}
	fmt.Println("Codeblock Language:", codeblock.Language)
	fmt.Println("Codeblock Content:", codeblock.Content)

	// Get the first argument
	argument := arguments.Get(0)

	// Parse it into an integer
	// HINT: You can also use the argument.AsInt64 method
	integer, err := argument.AsInt()
	if err != nil {
		// Argument is no integer
	}
	fmt.Println("Int:", integer)

	// Parse it into a boolean
	boolean, err := argument.AsBool()
	if err != nil {
		// Argument is no bool
	}
	fmt.Println("Bool:", boolean)

	// Parse it into a user ID if it is an user mention
	// HINT: You can also use the argument.AsRoleMentionID and argument.AsChannelMentionID methods
	userID := argument.AsUserMentionID()
	if userID == "" {
		// Argument is no user mention
	}
	fmt.Println("User ID:", userID)

	// Parse it into a duration
	dur, err := argument.AsDuration()
	if err != nil {
		// Argument is no duration
	}
	fmt.Println("Duration:", dur.String())
}
