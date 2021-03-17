package main

import "fmt"

type BalancingStrategy interface {
	Init([]*Backend)
	GetNextBackend(IncomingReq) *Backend
	RegisterBackend(*Backend)
	PrintTopology()
}

type RRBalancingStrategy struct {
	Index    int
	Backends []*Backend
}

type StaticBalancingStrategy struct {
	Index    int
	Backends []*Backend
}

func (s *RRBalancingStrategy) Init(backends []*Backend) {
	s.Index = 0
	s.Backends = backends
}

func (s *RRBalancingStrategy) GetNextBackend(_ IncomingReq) *Backend {
	s.Index = (s.Index + 1) % len(s.Backends)
	return s.Backends[s.Index]
}

func (s *RRBalancingStrategy) RegisterBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func (s *RRBalancingStrategy) PrintTopology() {
	for index, backend := range s.Backends {
		fmt.Println(fmt.Sprintf("      [%d] %s", index, backend))
	}
}

func NewRRBalancingStrategy(backends []*Backend) *RRBalancingStrategy {
	strategy := new(RRBalancingStrategy)
	strategy.Init(backends)
	return strategy
}

func (s *StaticBalancingStrategy) Init(backends []*Backend) {
	s.Index = 0
	s.Backends = backends
}

func (s *StaticBalancingStrategy) GetNextBackend(_ IncomingReq) *Backend {
	return s.Backends[s.Index]
}

func (s *StaticBalancingStrategy) RegisterBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func (s *StaticBalancingStrategy) PrintTopology() {
	for index, backend := range s.Backends {
		if index == s.Index {
			fmt.Println(fmt.Sprintf("      [%s] %s", "x", backend))
		} else {
			fmt.Println(fmt.Sprintf("      [%s] %s", " ", backend))
		}
	}
}

func NewStaticBalancingStrategy(backends []*Backend) *StaticBalancingStrategy {
	strategy := new(StaticBalancingStrategy)
	strategy.Init(backends)
	return strategy
}
