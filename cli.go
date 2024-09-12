package main

import (
	"fmt"
	"os"
)

func receiveCli() {
	args := os.Args[1:]
	cmdCount := len(args)

	switch {
	case cmdCount < 1:
		fmt.Println("no website provided")
		os.Exit(1)
	case cmdCount > 1:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	default:
		crawlHtml(args[0]) // Use the first command-line argument
	}

}
