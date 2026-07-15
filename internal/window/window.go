package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"clip/internal/clipboard"
	"clip/internal/config"
	"clip/internal/history"
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
			o.(*widget.Label).SetText(hist.At(i))
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
