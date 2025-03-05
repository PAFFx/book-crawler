package crawler

import (
	"fmt"
	"sync"
	"time"
)

// ProgressTracker tracks progress for a domain-specific crawler
type ProgressTracker struct {
	domain         string
	visitedCount   int
	extractedCount int
	lastURL        string
	startTime      time.Time
	mu             sync.Mutex
}

// ProgressManager manages multiple domain-specific progress trackers
type ProgressManager struct {
	trackers map[string]*ProgressTracker
	mutex    sync.Mutex
}

var (
	// Global progress manager instance
	globalManager     *ProgressManager
	globalManagerOnce sync.Once
)

// GetProgressManager returns the singleton instance of the progress manager
func GetProgressManager() *ProgressManager {
	globalManagerOnce.Do(func() {
		globalManager = NewProgressManager()
	})
	return globalManager
}

// NewProgressManager creates a new progress manager
func NewProgressManager() *ProgressManager {
	return &ProgressManager{
		trackers: make(map[string]*ProgressTracker),
	}
}

// GetTracker returns a domain-specific progress tracker
func (m *ProgressManager) GetTracker(domain string) *ProgressTracker {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if tracker, exists := m.trackers[domain]; exists {
		return tracker
	}

	// Create a new tracker for this domain
	tracker := &ProgressTracker{
		domain:    domain,
		startTime: time.Now(),
	}

	m.trackers[domain] = tracker
	return tracker
}

// LogProgress prints a summary of all progress trackers
func (m *ProgressManager) LogProgress() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fmt.Println("\n--- Crawler Progress Summary ---")
	for domain, tracker := range m.trackers {
		tracker.mu.Lock()
		elapsed := time.Since(tracker.startTime).Round(time.Second)
		fmt.Printf("Domain: %s - Visited: %d pages, Extracted: %d books, Running for: %s\n",
			domain, tracker.visitedCount, tracker.extractedCount, elapsed)
		if tracker.lastURL != "" {
			fmt.Printf("  Last URL: %s\n", tracker.lastURL)
		}
		tracker.mu.Unlock()
	}
	fmt.Println("--------------------------------")
}

// LogVisit records a page visit
func (t *ProgressTracker) LogVisit(url string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.visitedCount++
	t.lastURL = url
}

// LogExtraction records a book extraction
func (t *ProgressTracker) LogExtraction(url string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.extractedCount++
	t.lastURL = url
}

// StartPeriodicLogging starts a goroutine that logs progress at regular intervals
func (m *ProgressManager) StartPeriodicLogging(interval time.Duration) func() {
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.LogProgress()
			case <-stopCh:
				return
			}
		}
	}()

	return func() {
		close(stopCh)
	}
}
