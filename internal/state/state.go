package state

import (
	"sync"
	"time"
)

const (
	MaxLatency  = 30 * time.Second
	MaxJitter   = 10 * time.Second
	MaxDropRate = 1
)

type Chaos struct {
	Jitter   time.Duration
	Latency  time.Duration
	DropRate float64
}

type Connection struct {
	Port   int
	Target string
}

type HavocState struct {
	mu sync.RWMutex
	Chaos
	Connection
}

func New() *HavocState {
	return &HavocState{}
}

func (s *HavocState) SetPort(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Port = port
}

func (s *HavocState) SetTarget(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Target = target
}

func (s *HavocState) SetJitter(jitter time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jitter > MaxJitter {
		jitter = MaxJitter
	}

	if jitter < 0 {
		jitter = 0
	}

	s.Jitter = jitter
}

func (s *HavocState) SetLatency(latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if latency > MaxLatency {
		latency = MaxLatency
	}

	if latency < 0 {
		latency = 0
	}

	s.Latency = latency
}

func (s *HavocState) SetDropRate(dropRate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if dropRate > MaxDropRate {
		dropRate = MaxDropRate
	}

	if dropRate < 0 {
		dropRate = 0
	}

	s.DropRate = dropRate
}

func (s *HavocState) GetChaos() *Chaos {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &s.Chaos
}

func (s *HavocState) GetConnection() *Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &s.Connection
}
