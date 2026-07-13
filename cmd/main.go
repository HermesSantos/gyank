package main

import (
	"clip/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

var items []string

func main() {
	a := app.New()
	w := a.NewWindow("Clipboard")
	w.Resize(fyne.NewSize(400, 200))

	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(items[i])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		item := items[id]

		clipboard.Write(
			clipboard.FmtText,
			[]byte(item),
		)
	}

	w.SetContent(list)

	clipboardChan := internal.InitClipboard()

	go func() {
		for item := range clipboardChan {
			items = append(items, item)

			fyne.Do(func() {
				list.Refresh()
			})
		}
	}()

	w.ShowAndRun()
}
