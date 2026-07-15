package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Setup attaches a system tray menu. OS-specific tray details stay here.
func Setup(a fyne.App, show func(), quit func()) {
	desk, ok := a.(desktop.App)
	if !ok {
		return
	}

	desk.SetSystemTrayMenu(fyne.NewMenu(
		"",
		fyne.NewMenuItem("Show", show),
		fyne.NewMenuItem("Quit", quit),
	))
}
