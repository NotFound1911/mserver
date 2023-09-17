package memory

import (
	"context"
	"fmt"
	"github.com/NotFound1911/mserver/session"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Store struct {
	mu         sync.RWMutex
	c          *cache.Cache
	expiration time.Duration
}

var _ session.Store = &Store{}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		c:          cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (s *Store) Create(ctx context.Context, id string) (session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess := &memorySession{
		id:   id,
		data: map[string]string{},
	}
	s.c.Set(sess.ID(), sess, s.expiration)
	return sess, nil
}
func (s *Store) Update(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.c.Get(id)
	if !ok {
		return fmt.Errorf("session(%s) not found", id)
	}
	s.c.Set(sess.(*memorySession).ID(), sess, s.expiration)
	return nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.Delete(id)
	return nil
}
func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.c.Get(id)
	if !ok {
		return nil, fmt.Errorf("session(%s) not found", id)
	}
	return sess.(*memorySession), nil
}

// 基于内存存储
type memorySession struct {
	mu   sync.RWMutex
	id   string
	data map[string]string
}

var _ session.Session = &memorySession{}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key:%s is not exist", key)
	}
	return val, nil
}
func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
	return nil
}
func (m *memorySession) ID() string {
	return m.id
}
