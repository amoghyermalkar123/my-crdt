package pncounter

import "crdtsgo/gcounter"

type PNCounter struct {
	incCounter *gcounter.GCounter
	decCounter *gcounter.GCounter
}

func (p *PNCounter) Increment(val int64) int64 {
	return p.incCounter.Increment(val)
}

func (p *PNCounter) Decrement(val int64) int64 {
	return p.decCounter.Increment(val)
}

func (p *PNCounter) Value() int64 {
	return p.incCounter.Value() - p.decCounter.Value()
}
