package cache

import (
	"sync"
	"time"
)

const cleanChanBufSize = 100

type Cache interface {
	Add(key string, link string)
	Get(key string) (string, bool)
	Len() int
}

type cachedLink struct {
	key     string
	value   string
	timer   *time.Timer
	expires time.Time
}

type MemoryCache struct {
	MaxCap  int
	links   map[string]*cachedLink
	mu      sync.RWMutex
	cleaner chan string
	done    chan struct{}
	ttl     time.Duration
}

func NewMemoryCache(cap int, ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		MaxCap:  cap,
		ttl:     ttl,
		links:   make(map[string]*cachedLink, cap),
		cleaner: make(chan string, cleanChanBufSize),
		done:    make(chan struct{}),
	}
}

func (m *MemoryCache) Add(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if link, ok := m.links[key]; ok {
		m.updateLink(link, value)
		return
	}
	m.insertLink(key, value)
}

func (m *MemoryCache) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ent, ok := m.links[key]; ok {
		//проверяем если cleaner не успел почистить
		if m.ttl == 0 || time.Now().Before(ent.expires) {
			return ent.value, true
		}
	}

	return "", false
}
func (m *MemoryCache) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.links)
}

func (m *MemoryCache) RunCleaner() {

	for {
		select {
		case <-m.done:
			return
		case key := <-m.cleaner:
			m.DeleteKey(key)
		default:
			time.Sleep(100 * time.Microsecond) // чтобы не сжирать прецессное время
		}
	}

}

func (m *MemoryCache) DeleteKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.links, key)
}

func (m *MemoryCache) StopCleaner() {
	close(m.done)
}

func (m *MemoryCache) updateLink(e *cachedLink, value string) {
	// должен уже иметь блокировку записи
	e.value = value
	m.resetLinkTTL(e)
}

func (m *MemoryCache) resetLinkTTL(e *cachedLink) {
	// должен уже иметь блокировку записи
	if m.ttl > 0 {
		e.timer.Reset(m.ttl)
	}

	e.expires = time.Now().Add(m.ttl)
}
func (m *MemoryCache) insertLink(key, value string) *cachedLink {
	// должен уже иметь блокировку записи
	ent := &cachedLink{
		key:     key,
		value:   value,
		expires: time.Now().Add(m.ttl),
	}

	if len(m.links) >= m.MaxCap {
		return nil
	}

	if m.ttl > 0 {
		ent.timer = time.AfterFunc(m.ttl, func() {
			select {
			case <-m.done:
				return
			case m.cleaner <- key:
			}
		})
	}

	m.links[key] = ent
	return ent
}
