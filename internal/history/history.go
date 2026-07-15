package history

import "sync"

// History keeps an ordered, deduplicated list of clipboard entries.
// UI and clipboard packages must not own this logic.
type History struct {
	mu    sync.RWMutex
	items []string
}

func New() *History {
	return &History{}
}

// NewWith creates a history preloaded with items (e.g. from storage).
func NewWith(items []string) *History {
	copied := make([]string, len(items))
	copy(copied, items)
	return &History{items: copied}
}

// Add appends text if it is not already present. Returns true when added.
func (h *History) Add(text string) bool {
	if text == "" {
		return false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, item := range h.items {
		if item == text {
			return false
		}
	}

	h.items = append(h.items, text)
	return true
}

func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.items)
}

func (h *History) At(i int) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.items[i]
}

// Items returns a copy of the current entries.
func (h *History) Items() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]string, len(h.items))
	copy(out, h.items)
	return out
}
