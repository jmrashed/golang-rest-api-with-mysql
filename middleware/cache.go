package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Data      []byte
	Headers   http.Header
	Status    int
	ExpiresAt time.Time
}

// Cache holds cached responses
type Cache struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

// NewCache creates a new cache
func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Get retrieves cached entry
func (c *Cache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.entries[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	return entry, true
}

// Set stores cache entry
func (c *Cache) Set(key string, data []byte, headers http.Header, status int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries[key] = &CacheEntry{
		Data:      data,
		Headers:   headers,
		Status:    status,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// cleanup removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// CacheResponseWriter wraps http.ResponseWriter for caching
type CacheResponseWriter struct {
	http.ResponseWriter
	data   []byte
	status int
}

func NewCacheResponseWriter(w http.ResponseWriter) *CacheResponseWriter {
	return &CacheResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (crw *CacheResponseWriter) WriteHeader(status int) {
	crw.status = status
	crw.ResponseWriter.WriteHeader(status)
}

func (crw *CacheResponseWriter) Write(data []byte) (int, error) {
	crw.data = append(crw.data, data...)
	return crw.ResponseWriter.Write(data)
}

// CacheMiddleware creates caching middleware
func CacheMiddleware(cache *Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}
			
			// Generate cache key
			key := generateCacheKey(r)
			
			// Check cache
			if entry, exists := cache.Get(key); exists {
				// Copy headers
				for k, v := range entry.Headers {
					w.Header()[k] = v
				}
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(entry.Status)
				w.Write(entry.Data)
				return
			}
			
			// Cache miss - process request
			crw := NewCacheResponseWriter(w)
			next.ServeHTTP(crw, r)
			
			// Cache successful responses
			if crw.status == http.StatusOK {
				cache.Set(key, crw.data, crw.Header(), crw.status)
				w.Header().Set("X-Cache", "MISS")
			}
		})
	}
}

// generateCacheKey creates a cache key from request
func generateCacheKey(r *http.Request) string {
	key := fmt.Sprintf("%s:%s:%s", r.Method, r.URL.Path, r.URL.RawQuery)
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}