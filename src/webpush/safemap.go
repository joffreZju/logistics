package main

import "sync"

type PushPool struct {
	S  map[int]*chan []byte
	mu sync.RWMutex
}

func NewPushPool() *PushPool {
	p := new(PushPool)
	p.S = make(map[int]*chan []byte)
	return p
}

func (p *PushPool) AddUser(id int, ch *chan []byte) {
	p.mu.Lock()
	p.S[id] = ch
	p.mu.Unlock()
}

func (p *PushPool) GetUserChan(id int) (ch *chan []byte, b bool) {
	p.mu.RLock()
	ch, b = p.S[id]
	p.mu.RUnlock()
	return
}
