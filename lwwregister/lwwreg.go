package lwwregister

import "crdtsgo/gcounter"

const CLOCK_INCR = 1

type LWWRegister struct {
	register map[string]*Value
	nodeID   string
}

type Value struct {
	baseClock *gcounter.GCounter
	value     string
}

func (l *LWWRegister) Update(key, value string) {
	if val, ok := l.register[key]; ok {
		currentCounter := val.baseClock.Increment(CLOCK_INCR)
		if currentCounter > val.baseClock.Value() {
			val.value = value
			val.baseClock.Increment(currentCounter)
		}
	}
}

func (l *LWWRegister) GetValueByKey(key string) string {
	if value, ok := l.register[key]; ok {
		return value.value
	}
	return ""
}
