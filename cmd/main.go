package main

import (
	"fmt"
	"time"

	"golang.design/x/clipboard"
)

func main() {
	clipboard.Init()

	var last string
	var lastOnes []string

	for {
		current := string(clipboard.Read(clipboard.FmtText))

		if current != last {
			fmt.Println("Novo conteúdo:")
			fmt.Println(current)

			last = current
			lastOnes = append(lastOnes, current)

			fmt.Println("Lista:")
			fmt.Println(lastOnes)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
