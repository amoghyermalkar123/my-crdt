package vectorclocks

import (
	"crdtsgo/gcounter"
	"sync"
)

const CLOCK_INCREMENT_OFFSET = 1

type VectorClock struct {
	clockMap    map[string]BaseClock
	localNodeID string
	mu          sync.RWMutex
}

type BaseClock struct {
	clockID *gcounter.GCounter
}

func (vc *VectorClock) IncrementLocalMonotonicClock() {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	if baseClock, ok := vc.clockMap[vc.localNodeID]; ok {
		baseClock.clockID.Increment(CLOCK_INCREMENT_OFFSET)
	}
}

func (vc *VectorClock) IntegrateClocks(integrationClock map[string]BaseClock) {
	for nodeID, remoteReplicaClock := range integrationClock {
		if nodeID == vc.localNodeID {
			continue
		}
		if localReplicaClock, ok := vc.clockMap[nodeID]; ok {
			localReplicaClock.clockID.Merge(remoteReplicaClock.clockID.Value())
		} else {
			vc.clockMap[nodeID] = remoteReplicaClock
		}
	}
}
