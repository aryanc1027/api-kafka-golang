package models

import (
	"sync"
	"time"
)

type Stream struct {
	ID        string
	CreatedAt time.Time
	Messages  chan []byte
	Errors    chan error
	mu        sync.Mutex
	closed    bool
}

func NewStream(id string) *Stream {
	return &Stream{
		ID:        id,
		CreatedAt: time.Now(),
		Messages:  make(chan []byte, 100), // Buffered channel
		Errors:    make(chan error, 10),   // Buffered channel
	}
}

func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		close(s.Messages)
		close(s.Errors)
		s.closed = true
	}
}

func (s *Stream) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}