package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func printBackends() {
	fmt.Println("Backends:")
	backends := lb.backends
	if len(backends) == 0 {
		fmt.Println("  no backends")
	} else {
		for _, backend := range backends {
			fmt.Println(fmt.Sprintf("  - %s", backend))
		}
	}
}

func cli() {
	for {
		var command string
		fmt.Print(">>> ")
		fmt.Scanf("%s", &command)
		tokens := strings.Split(command, " ")
		switch tokens[0] {
		case "quit":
			lb.events <- "quit"
			// TODO: this is not idea. End this gracefully.
			return
		case "exit":
			lb.events <- "quit"
			// TODO: this is not idea. End this gracefully.
			return
		case "list":
			printBackends()
		}
	}
}

func setLogging() {
	f, err := os.OpenFile("lb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	// TODO: Closse the log file.
}

func main() {
	setLogging()
	go lb.Run()
	cli()
}
