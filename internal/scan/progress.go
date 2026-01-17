package scan

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type ProgressTracker struct {
	total     int64
	processed int64
	mutex     sync.Mutex
	writer    io.Writer
	enabled   bool
	done      chan struct{}
	ticker    *time.Ticker
}

func NewProgressTracker(total int64, enabled bool) *ProgressTracker {
	return &ProgressTracker{
		total:   total,
		writer:  os.Stderr,
		enabled: enabled,
		done:    make(chan struct{}),
	}
}

func (pt *ProgressTracker) Start() {
	if !pt.enabled || pt.total == 0 {
		return
	}

	pt.ticker = time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-pt.ticker.C:
				pt.render()
			case <-pt.done:
				pt.render()
				fmt.Fprintln(pt.writer, "\n")
				return
			}
		}
	}()
}

func (pt *ProgressTracker) Increment() {
	if !pt.enabled {
		return
	}

	pt.mutex.Lock()
	pt.processed++
	pt.mutex.Unlock()
}

func (pt *ProgressTracker) render() {
	pt.mutex.Lock()
	processed := pt.processed
	total := pt.total
	pt.mutex.Unlock()

	if total == 0 {
		return
	}

	percentage := float64(processed) / float64(total) * 100
	barWidth := 40
	filled := int(float64(barWidth) * float64(processed) / float64(total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Fprintf(pt.writer, "\r[%s] %3.0f%% (%d/%d files)",
		bar, percentage, processed, total)
}

func (pt *ProgressTracker) Stop() {
	if !pt.enabled {
		return
	}
	if pt.ticker != nil {
		pt.ticker.Stop()
	}
	close(pt.done)
}
