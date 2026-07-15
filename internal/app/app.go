package app

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"

	"clip/internal/clipboard"
	"clip/internal/config"
	"clip/internal/history"
	"clip/internal/tray"
	"clip/internal/window"
)

// Run wires dependencies and starts the long-running desktop process.
func Run() error {
	cfg := config.Load()

	clip, err := clipboard.New()
	if err != nil {
		return fmt.Errorf("clipboard: %w", err)
	}

	hist := history.New()
	fyneApp := fyneapp.New()
	histWin := window.NewHistory(fyneApp, hist, clip, cfg)

	tray.Setup(fyneApp, histWin.Show, fyneApp.Quit)

	ctx, cancel := context.WithCancel(context.Background())
	fyneApp.Lifecycle().SetOnStopped(cancel)

	go watchClipboard(ctx, clip, hist, histWin, cfg)

	histWin.Window().ShowAndRun()
	return nil
}

func watchClipboard(
	ctx context.Context,
	clip *clipboard.Clipboard,
	hist *history.History,
	win *window.History,
	cfg config.Config,
) {
	for text := range clip.Watch(ctx, cfg.PollInterval) {
		if !hist.Add(text) {
			continue
		}
		fyne.Do(func() {
			win.Refresh()
		})
	}
	log.Println("clipboard watch stopped")
}
