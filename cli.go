package main

import (
	"fmt"
	"os"
	"strconv"
)

func receiveCli() {
	args := os.Args[1:]
	cmdCount := len(args)

	switch {
	case cmdCount < 3:
		fmt.Println("no website provided")
		os.Exit(1)
	case cmdCount > 3:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	default:
		maxConcurrent, _ := strconv.Atoi(args[1])
		maxPages, _ := strconv.Atoi(args[2])
		crawlHtml(args[0], maxConcurrent, maxPages) // Use the first command-line argument
	}

}
