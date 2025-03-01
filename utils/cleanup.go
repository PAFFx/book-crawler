package utils

import (
	"fmt"
	"sync"
)

// CleanupManager struct to manage multiple cleanup functions
type CleanupManager struct {
	mu       sync.Mutex
	handlers []func()
}

// Global cleanup manager instance
var instance *CleanupManager
var once sync.Once

// GetCleanupManager returns a singleton instance of CleanupManager
func GetCleanupManager() *CleanupManager {
	once.Do(func() {
		instance = &CleanupManager{}
	})
	return instance
}

// Add registers a new cleanup function
func (c *CleanupManager) Add(fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers = append(c.handlers, fn)
}

// RunAll executes all cleanup functions in reverse order before exiting
func (c *CleanupManager) RunAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("\nExecuting all cleanup functions before exit...")

	// Run cleanup functions in reverse order (like defer)
	for i := len(c.handlers) - 1; i >= 0; i-- {
		c.handlers[i]()
	}
}
