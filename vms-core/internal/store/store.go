package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Store interface {
	Save(key string, value any) error
	Load(key string, value any) error
}

var _ Store = (*FileStore)(nil)

type FileStore struct {
	path string
	mu   sync.RWMutex
	data map[string]any
}

func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{
		path: path,
		data: make(map[string]any),
	}
	if err := fs.loadFromFile(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStore) loadFromFile() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.ReadFile(fs.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, start with empty data
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, &fs.data)
}

func (fs *FileStore) saveToFile() error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := json.MarshalIndent(fs.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.path, data, 0644)
}

func (fs *FileStore) Save(key string, value any) error {
	fs.mu.Lock()
	fs.data[key] = value
	fs.mu.Unlock()

	return fs.saveToFile()
}

func (fs *FileStore) Load(key string, value any) error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	val, ok := fs.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	// Hack to handle JSON unmarshalling into the target pointer
	// Since we are effectively working with interface{}, we marshal back to json
	// and unmarshal into the target to handle type conversion robustly (like map to struct)
	// This is a simple but inefficient way. For simple types it works.
	// For this specific use case (string), it should be fine.

	// Optimization for string
	if vStr, ok := val.(string); ok {
		if destStr, ok := value.(*string); ok {
			*destStr = vStr
			return nil
		}
	}

	// Fallback generic approach
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
