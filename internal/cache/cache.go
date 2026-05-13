package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const cacheTTL = 72 * time.Hour

type CacheEntry struct {
	QuoteID   string    `json:"quote_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Classify  string    `json:"classify"`
	CachedAt  time.Time `json:"cached_at"`
}

var (
	store sync.Map
	mu    sync.Mutex
)

func cachePath() string {
	dir := os.Getenv("EFINANCE_CACHE_DIR")
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".efinance")
	}
	return filepath.Join(dir, "search-cache.json")
}

func Load() error {
	p := cachePath()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var m map[string]CacheEntry
	if err := json.Unmarshal(data, &m); err != nil {
		_ = os.WriteFile(p, []byte{}, 0644)
		return nil
	}
	for k, v := range m {
		store.Store(k, v)
	}
	return nil
}

func Get(code string) (CacheEntry, bool) {
	v, ok := store.Load(code)
	if !ok {
		return CacheEntry{}, false
	}
	entry := v.(CacheEntry)
	if time.Since(entry.CachedAt) > cacheTTL {
		store.Delete(code)
		return CacheEntry{}, false
	}
	return entry, true
}

func Set(code string, entry CacheEntry) {
	entry.CachedAt = time.Now()
	store.Store(code, entry)
	go save()
}

func save() {
	mu.Lock()
	defer mu.Unlock()
	m := make(map[string]CacheEntry)
	store.Range(func(key, value any) bool {
		e := value.(CacheEntry)
		if time.Since(e.CachedAt) <= cacheTTL {
			m[key.(string)] = e
		}
		return true
	})
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	p := cachePath()
	os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, data, 0644)
}
