package main

import (
	"github.com/mattn/go-runewidth"

	"github.com/nsf/termbox-go"
)

type (
	GameMode uint8
)

const (
	GameModeStop GameMode = iota
	GameModeStart
)

type Game struct {
	mode      GameMode
	questions []string
	inputs    []int32
}

func NewGame() *Game {
	g := Game{
		// TODO: intiialize random words / random length
		questions: []string{
			"the",
			"typing",
			"game",
			"very",
			"very",
			"hard",
		},
	}
	return &g
}

func (g *Game) start() {
	g.mode = GameModeStart
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
}
func (g *Game) stop() {
	g.mode = GameModeStop
}

func (g *Game) Update(eventQueue chan termbox.Event) bool {
	ev := <-eventQueue

	// Escapeでゲーム終了
	if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
		return false
	}

	switch g.mode {
	case GameModeStop:
		// Enterでゲーム開始
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
			g.start()
		}
	case GameModeStart:
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyBackspace, termbox.KeyDelete, termbox.KeyBackspace2:
				if len(g.inputs) > 0 {
					g.inputs = g.inputs[:len(g.inputs)-1]
				}
			default:
				if ev.Ch >= 'a' && ev.Ch <= 'z' {
					g.inputs = append(g.inputs, ev.Ch)
				}
			}
		}
	}

	return true
}

func (g *Game) Render() {
	switch g.mode {
	case GameModeStart:
		g.renderStart()
	case GameModeStop:
		g.renderStop()
	}
}

func (g *Game) renderStart() {
	var current int
	for i, q := range g.questions {
		for j, r := range q {
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault
			if current < len(g.inputs) {
				if g.inputs[current] == r {
					fg = termbox.ColorGreen
				} else {
					fg = termbox.ColorYellow
				}
			}
			termbox.SetCell(j, i, r, fg, bg)
			current++
		}
	}
}

func (g *Game) renderStop() {
	w, h := termbox.Size()

	messages := [...]string{
		"-- Typing Game --",
		"",
		"Start: Press ENTER",
		"Exit: ESC",
	}
	dy := (h - len(messages)) / 2

	for y, message := range messages {
		dx := (w - runewidth.StringWidth(message)) / 2
		for x, r := range message {
			termbox.SetCell(dx+x, dy+y, r, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
}

func mainLoop() {
	g := NewGame()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	render := func(g *Game) error {
		g.Render()
		return termbox.Flush()
	}

	if err := render(g); err != nil {
		return
	}
	for {
		if !g.Update(eventQueue) {
			return
		}
		if err := render(g); err != nil {
			return
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	mainLoop()
}
