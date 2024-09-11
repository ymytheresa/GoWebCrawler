package main

import (
	"bufio"
	"os"
	"strings"
)

func receiveCli(){
	reader := bufio.NewReader(os.Stdin)
	for {
		command, _ := reader. ReadString('\n')
		commandSlice := strings.Fields(command)
		switch len(commandSlice)
	}
}