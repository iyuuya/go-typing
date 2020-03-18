package main

import (
	"github.com/iyuuya/go-typing/internal/typing"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	g := typing.NewGame()
	g.Go()
}
