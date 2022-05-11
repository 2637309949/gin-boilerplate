package async

import (
	"sync/atomic"
)

const (
	start int32 = iota
	running
	shutdown
	stop
)

type simpleFSM struct {
	status int32
}

func newSimpleFSM() *simpleFSM {
	return &simpleFSM{
		status: start,
	}
}

func (s *simpleFSM) actEvent(stat int32) bool {
	if s.Current() > stat {
		return false
	} else {
		return atomic.CompareAndSwapInt32(&s.status, s.status, stat)
	}
}

func (s *simpleFSM) isRunning() bool {
	return s.Current() == running
}

func (s *simpleFSM) Current() int32 {
	return atomic.LoadInt32(&s.status)
}

func (s *simpleFSM) Action(stat int32) bool {
	switch stat {
	case start:
	case running:
	case shutdown:
	case stop:
		return s.actEvent(stat)
	default:
		return false
	}
	return false
}
