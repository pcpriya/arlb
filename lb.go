package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Backend struct {
	Host      string
	Port      int
	IsHealthy bool
}

func (b Backend) String() string {
	return fmt.Sprintf("%s:%d healthy:%v", b.Host, b.Port, b.IsHealthy)
}

type LB struct {
	backends []Backend
	events   chan string
}

var lb LB

func init() {
	lb = LB{}
	lb.events = make(chan string)
	lb.backends = []Backend{
		Backend{Host: "localhost", Port: 8081, IsHealthy: false},
		Backend{Host: "localhost", Port: 8082, IsHealthy: true},
	}
}

func (l LB) Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	log.Println("LB listening on port 8080 ...")

	go func() {
		for {
			select {
			case event := <-l.events:
				if event == "quit" {
					log.Println("gracefully terminating ...")
					return
				}
			}
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		// Once the connection is accepted proxying it to backend
		go lb.proxy(connection)
	}
}

func (l LB) proxy(srcConnection net.Conn) {
	// Get backend sserver depending on some algorithm
	backend := l.backends[0]

	// Setup backend connection
	backendServerConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
	if err != nil {
		log.Printf("error connecting to backend. %s", err.Error())

		// send back error to src
		srcConnection.Write([]byte("backend not available"))
		srcConnection.Close()
	}

	go io.Copy(backendServerConnection, srcConnection)
	go io.Copy(srcConnection, backendServerConnection)
}
