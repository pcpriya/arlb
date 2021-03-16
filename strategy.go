package main

type BalancingStrategy interface {
	GetNextBackend([]*Backend) *Backend
}

type RoundRobinBalancingStrategy struct {
	Index int
}

func (s *RoundRobinBalancingStrategy) GetNextBackend(backends []*Backend) *Backend {
	s.Index = (s.Index + 1) % len(backends)
	return backends[s.Index]
}

var STRATEGY_ROUNDROBIN *RoundRobinBalancingStrategy

func InitStrategy() {
	STRATEGY_ROUNDROBIN = &RoundRobinBalancingStrategy{Index: 0}
}
