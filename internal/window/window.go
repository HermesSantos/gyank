package window

import (
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"clip/internal/clipboard"
	"clip/internal/config"
	"clip/internal/history"
)

const (
	previewMaxWords = 8
	previewMaxRunes = 80
)

// History is the main window that displays clipboard history.
// It only renders state and forwards user actions; no business rules live here.
type History struct {
	win  fyne.Window
	list *widget.List
	hist *history.History
}

func NewHistory(a fyne.App, hist *history.History, clip *clipboard.Clipboard, cfg config.Config) *History {
	w := a.NewWindow(cfg.AppName)
	w.Resize(fyne.NewSize(cfg.WindowWidth, cfg.WindowHeight))

	h := &History{win: w, hist: hist}

	h.list = widget.NewList(
		func() int {
			return hist.Len()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(preview(hist.At(i)))
		},
	)

	h.list.OnSelected = func(id widget.ListItemID) {
		clip.Write(hist.At(id))
	}

	w.SetContent(h.list)
	return h
}

func (h *History) Window() fyne.Window {
	return h.win
}

func (h *History) Show() {
	h.win.Show()
}

func (h *History) Refresh() {
	h.list.Refresh()
}

// preview shortens text for the list label. Full text stays in history.
func preview(text string) string {
	words := strings.Fields(text)
	truncated := false

	if len(words) > previewMaxWords {
		words = words[:previewMaxWords]
		truncated = true
	}

	out := strings.Join(words, " ")
	if utf8.RuneCountInString(out) > previewMaxRunes {
		runes := []rune(out)
		out = string(runes[:previewMaxRunes])
		truncated = true
	}

	if truncated {
		return out + "..."
	}
	return out
}
