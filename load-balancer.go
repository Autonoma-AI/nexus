package main

import (
	"sync"
)

type LoadBalancer struct {
	currentServiceIndex map[string]int
	mutex               sync.Mutex
}

func (lb *LoadBalancer) getNextBackend(service string, backends []*Backend) (*Backend, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if lb.currentServiceIndex == nil {
		lb.currentServiceIndex = make(map[string]int)
	}

	index := lb.currentServiceIndex[service]
	backend := backends[index]
	lb.currentServiceIndex[service] = (lb.currentServiceIndex[service] + 1) % len(backends)
	return backend, nil
}
