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
		Backend{Host: "localhost", Port: 8083, IsHealthy: true},
		Backend{Host: "localhost", Port: 8084, IsHealthy: true},
	}
}

func (l LB) Run() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	log.Println("LB listening on port 9090 ...")

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
			panic(err)
		}

		log.Printf(connection.LocalAddr().String(), connection.RemoteAddr().String())

		// Once the connection is accepted proxying it to backend
		go lb.proxy(connection)
	}
}

var index int = 0

func (l LB) proxy(srcConnection net.Conn) {
	index = (index + 1) % len(l.backends)

	// Get backend sserver depending on some algorithm
	backend := l.backends[index]

	// Setup backend connection
	backendServerConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
	if err != nil {
		log.Printf("error connecting to backend. %s", err.Error())

		// send back error to src
		srcConnection.Write([]byte("backend not available"))
		srcConnection.Close()
		panic(err)
	}

	go io.Copy(backendServerConnection, srcConnection)
	go io.Copy(srcConnection, backendServerConnection)
}
