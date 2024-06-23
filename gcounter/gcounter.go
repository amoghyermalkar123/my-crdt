package gcounter

import "sync"

type GCounter struct {
	counterValue int64
	mu           sync.RWMutex
}

func New() *GCounter {
	return &GCounter{counterValue: 0}
}

func (g *GCounter) Value() int64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.counterValue
}

func (g *GCounter) Increment(val int64) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counterValue = g.counterValue + val
	return g.counterValue
}

func (g *GCounter) Merge(integrationValue int64) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if integrationValue > g.counterValue {
		g.counterValue = integrationValue
	}
	return g.counterValue
}
