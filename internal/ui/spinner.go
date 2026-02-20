package ui

import (
	"fmt"
	"sync"
	"time"
)

// Spinner displays an animated loading indicator in the terminal.
type Spinner struct {
	done chan struct{}
	wg   sync.WaitGroup
}

// NewSpinner starts a spinner with the given message.
func NewSpinner(msg string) *Spinner {
	s := &Spinner{done: make(chan struct{})}
	s.wg.Add(1)
	go s.run(msg)
	return s
}

func (s *Spinner) run(msg string) {
	defer s.wg.Done()
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-s.done:
			fmt.Print("\r\033[K") // clear the spinner line
			return
		case <-ticker.C:
			fmt.Printf("\r  %s %s", DimStyle.Render(frames[i%len(frames)]), DimStyle.Render(msg))
			i++
		}
	}
}

// Stop halts the spinner and clears its line.
func (s *Spinner) Stop() {
	close(s.done)
	s.wg.Wait()
}
