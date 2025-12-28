package main

import "sync"

type wait struct {
	m  map[string]bool
	mu sync.Mutex
}

func newWait() *wait {
	return &wait{
		m: make(map[string]bool),
	}
}

func (w *wait) Add(s string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.m[s] = true
}

func (w *wait) Take(s string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.m[s] {
		return false
	}

	delete(w.m, s)
	return true
}
