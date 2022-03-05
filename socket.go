package main

import (
	"encoding/json"
	"sync"
)

type SocketStore struct {
	mu      sync.Mutex
	clients []*SocketConn
}

func (ss *SocketStore) PushStatus(s *status) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	res, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	for i := len(ss.clients) - 1; i >= 0; i-- {
		client := ss.clients[i]
		if client.IsClosed() {
			ss.clients = append(ss.clients[:i], ss.clients[i+1:]...)
			continue
		}
		client.Send(res)
	}
}

func (ss *SocketStore) AddSocketConn(sc *SocketConn) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.clients = append(ss.clients, sc)
}

type SocketConn struct {
	mu     sync.Mutex
	closed bool
	C      chan []byte
}

func NewSocketConn() *SocketConn {
	return &SocketConn{
		C:      make(chan []byte),
		closed: false,
	}
}

func (sc *SocketConn) SafeClose() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if !sc.closed {
		close(sc.C)
		sc.closed = true
	}
}

func (sc *SocketConn) IsClosed() bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.closed
}

func (sc *SocketConn) Send(data []byte) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return
	}

	sc.C <- data
}
