package main

import (
	"fmt"
	"log"
	"os"
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
		switch command {
		case "quit":
			lb.events <- Event{EventName: "quit"}
			// TODO: this is not idea. End this gracefully.
			return
		case "exit":
			lb.events <- Event{EventName: "quit"}
			// TODO: this is not idea. End this gracefully.
			return
		case "list":
			printBackends()
		case "add":
			var host string
			var port int

			fmt.Print("       Host: ")
			fmt.Scanf("%s", &host)

			fmt.Print("       Port: ")
			fmt.Scanf("%d", &port)

			lb.events <- Event{EventName: "backend/add", Data: Backend{Host: host, Port: port}}
		}
	}
}

func setLogging() {
	f, err := os.OpenFile("lb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		panic(err)
	}
	log.SetOutput(f)

	// TODO: Closse the log file.
}

func main() {
	setLogging()
	InitStrategy()
	InitLB()
	go lb.Run()
	cli()
}
