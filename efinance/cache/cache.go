package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache 通用缓存
type Cache struct {
	mu      sync.RWMutex
	data    map[string]cacheItem
	maxAge  time.Duration
	baseDir string
}

type cacheItem struct {
	Value     interface{} `json:"v"`
	Timestamp time.Time   `json:"t"`
}

// NewCache 创建缓存实例
func NewCache(maxAge time.Duration) *Cache {
	return &Cache{
		data:   make(map[string]cacheItem),
		maxAge: maxAge,
	}
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		Value:     value,
		Timestamp: time.Now(),
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if c.maxAge > 0 && time.Since(item.Timestamp) > c.maxAge {
		return nil, false
	}
	return item.Value, true
}

// Delete 删除缓存
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]cacheItem)
}

// Size 获取缓存大小
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// FileCache 基于文件的缓存
type FileCache struct {
	baseDir string
	maxAge  time.Duration
	mu      sync.RWMutex
}

// NewFileCache 创建文件缓存
func NewFileCache(baseDir string, maxAge time.Duration) (*FileCache, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &FileCache{
		baseDir: baseDir,
		maxAge:  maxAge,
	}, nil
}

// Set 保存到文件
func (c *FileCache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	filename := c.getFilename(key)
	return os.WriteFile(filename, data, 0644)
}

// Get 从文件读取
func (c *FileCache) Get(key string, dest interface{}) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filename := c.getFilename(key)

	info, err := os.Stat(filename)
	if err != nil {
		return false, nil
	}

	// 检查过期
	if c.maxAge > 0 && time.Since(info.ModTime()) > c.maxAge {
		return false, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}

	return true, nil
}

// Delete 删除文件
func (c *FileCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	filename := c.getFilename(key)
	os.Remove(filename)
	return nil
}

// Clear 清空所有缓存文件
func (c *FileCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		os.Remove(filepath.Join(c.baseDir, entry.Name()))
	}

	return nil
}

func (c *FileCache) getFilename(key string) string {
	return filepath.Join(c.baseDir, key+".json")
}
