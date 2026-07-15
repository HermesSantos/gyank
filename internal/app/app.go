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
	"clip/internal/storage"
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

	store, err := storage.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}

	items, err := store.Load()
	if err != nil {
		_ = store.Close()
		return fmt.Errorf("load history: %w", err)
	}

	hist := history.NewWith(items)
	fyneApp := fyneapp.New()
	histWin := window.NewHistory(fyneApp, hist, clip, cfg)

	tray.Setup(fyneApp, histWin.Show, fyneApp.Quit)

	ctx, cancel := context.WithCancel(context.Background())
	fyneApp.Lifecycle().SetOnStopped(func() {
		cancel()
		if err := store.Close(); err != nil {
			log.Printf("close storage: %v", err)
		}
	})

	go watchClipboard(ctx, clip, hist, store, histWin, cfg)

	histWin.Window().ShowAndRun()
	return nil
}

func watchClipboard(
	ctx context.Context,
	clip *clipboard.Clipboard,
	hist *history.History,
	store *storage.Store,
	win *window.History,
	cfg config.Config,
) {
	for text := range clip.Watch(ctx, cfg.PollInterval) {
		if !hist.Add(text) {
			continue
		}
		if err := store.Insert(text); err != nil {
			log.Printf("persist entry: %v", err)
		}
		fyne.Do(func() {
			win.Refresh()
		})
	}
	log.Println("clipboard watch stopped")
}
