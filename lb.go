package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type Backend struct {
	Host        string
	Port        int
	IsHealthy   bool
	NumRequests int
}

func (b *Backend) String() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

type Event struct {
	EventName string
	Data      interface{}
}

type LB struct {
	backends []*Backend
	events   chan Event
	strategy BalancingStrategy
}

type IncomingReq struct {
	srcConn net.Conn
	reqId   string
}

var lb *LB

func InitLB() {
	lb = &LB{
		events: make(chan Event),
		backends: []*Backend{
			&Backend{Host: "localhost", Port: 8081, IsHealthy: true},
		},
		strategy: STRATEGY_ROUNDROBIN,
	}
}

func (lb *LB) Run() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	log.Println("LB listening on port 9090 ...")

	go func() {
		for {
			select {
			case event := <-lb.events:
				if event.EventName == "quit" {
					log.Println("gracefully terminating ...")
					return
				} else if event.EventName == "backend/add" {
					backend, isOk := event.Data.(Backend)
					if !isOk {
						panic(err)
					}
					lb.backends = append(lb.backends, &backend)
				} else if event.EventName == "strategy/change" {
					strategyName, isOk := event.Data.(string)
					if !isOk {
						panic(err)
					}
					switch strategyName {
					case "round-robin":
						lb.strategy = STRATEGY_ROUNDROBIN
					default:
						lb.strategy = STRATEGY_ROUNDROBIN
					}
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

		// log.Printf("local: %v remote: %v", connection.LocalAddr().String(), connection.RemoteAddr().String())

		// Once the connection is accepted proxying it to backend
		go lb.proxy(IncomingReq{
			srcConn: connection,
			reqId:   uuid.NewString(),
		})
	}
}

func (lb *LB) proxy(req IncomingReq) {
	// Get backend sserver depending on some algorithm
	backend := lb.strategy.GetNextBackend(lb.backends)
	log.Printf("in-req: %s out-req: %s", req.reqId, backend.String())

	// Setup backend connection
	backendServerConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
	if err != nil {
		log.Printf("error connecting to backend. %s", err.Error())

		// send back error to src
		req.srcConn.Write([]byte("backend not available"))
		req.srcConn.Close()
		panic(err)
	}

	backend.NumRequests++
	go io.Copy(backendServerConnection, req.srcConn)
	go io.Copy(req.srcConn, backendServerConnection)
}
