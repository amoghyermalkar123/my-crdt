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

// ponder: at it's core isn't a vector clock also a LWW register?
// we're comparing two maps of arbitrary types but the core operation
// remains the same which is the Merge() call where we check the timestamps
// i.e. baseClock value and keep the the clock with a higher value.
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
