package internal

import (
	"time"

	"golang.design/x/clipboard"
)

func InitClipboard() <-chan string {
	clipboard.Init()

	ch := make(chan string)

	go func() {
		var last string

		for {
			current := string(clipboard.Read(clipboard.FmtText))

			if current != last {
				last = current
				ch <- current
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return ch
}
