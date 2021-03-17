package main

type BalancingStrategy interface {
	Init([]*Backend)
	GetNextBackend() *Backend
	RegisterBackend(*Backend)
}

type RRBalancingStrategy struct {
	Index    int
	Backends []*Backend
}

func (s *RRBalancingStrategy) Init(backends []*Backend) {
	s.Index = 0
	s.Backends = backends
}

func (s *RRBalancingStrategy) GetNextBackend() *Backend {
	s.Index = (s.Index + 1) % len(s.Backends)
	return s.Backends[s.Index]
}

func (s *RRBalancingStrategy) RegisterBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func NewRRBalancingStrategy(backends []*Backend) *RRBalancingStrategy {
	strategy := new(RRBalancingStrategy)
	strategy.Init(backends)
	return strategy
}
