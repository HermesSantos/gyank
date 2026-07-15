package clipboard

import (
	"context"
	"fmt"
	"time"

	sysclip "golang.design/x/clipboard"
)

// Clipboard reads, writes, and watches the system clipboard.
// Isolated here so OS clipboard details stay out of UI and history.
type Clipboard struct{}

func New() (*Clipboard, error) {
	if err := sysclip.Init(); err != nil {
		return nil, fmt.Errorf("clipboard init: %w", err)
	}
	return &Clipboard{}, nil
}

func (c *Clipboard) Read() string {
	return string(sysclip.Read(sysclip.FmtText))
}

func (c *Clipboard) Write(text string) {
	sysclip.Write(sysclip.FmtText, []byte(text))
}

// Watch polls the clipboard and sends new text values until ctx is cancelled.
// The returned channel is closed when the watch stops.
func (c *Clipboard) Watch(ctx context.Context, interval time.Duration) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		var last string
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current := c.Read()
				if current == last {
					continue
				}
				last = current
				select {
				case ch <- current:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch
}
