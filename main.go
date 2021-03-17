package main

import (
	"fmt"
	"log"
	"os"
)

func printBackends() {
	lb.strategy.PrintTopology()
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
		case "add":
			var host string
			var port int

			fmt.Print("       Host: ")
			fmt.Scanf("%s", &host)

			fmt.Print("       Port: ")
			fmt.Scanf("%d", &port)

			lb.events <- Event{EventName: "backend/add", Data: Backend{Host: host, Port: port}}
		case "strategy":
			var strategy string

			fmt.Print("       Name of the strategy: ")
			fmt.Scanf("%s", &strategy)

			lb.events <- Event{EventName: "strategy/change", Data: strategy}
		case "test":
			var reqId string

			fmt.Print("       Request ID: ")
			fmt.Scanf("%s", &reqId)

			backend := lb.strategy.GetNextBackend(IncomingReq{reqId: reqId})
			fmt.Printf("request: %s goes to backend: %s\n", reqId, backend)
		case "topology":
			printBackends()
		default:
			fmt.Println("available commands: topology, add, strategy, exit")
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
	InitLB()
	go lb.Run()
	cli()
}
