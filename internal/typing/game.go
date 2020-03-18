package typing

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type gameMode uint8

const (
	gameModeStop gameMode = iota
	gameModeStart
	gameModeFinished
)

type game struct {
	mode      gameMode
	questions []string
	inputs    []int32
}

func NewGame() *game {
	g := game{
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

func (g *game) start() {
	g.mode = gameModeStart
	g.inputs = g.inputs[:0]
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
}

func (g *game) stop() {
	g.mode = gameModeStop
}

func (g *game) finish() {
	g.mode = gameModeFinished

	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	// TODO: 時間はかって出したり正解率出したり
}

func (g *game) update(eventQueue chan termbox.Event) bool {
	ev := <-eventQueue

	// Escapeでゲーム終了
	if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
		return false
	}
	if ev.Type == termbox.EventResize {
		return termbox.Clear(termbox.ColorWhite, termbox.ColorDefault) == nil
	}

	switch g.mode {
	case gameModeStop, gameModeFinished:
		// Enterでゲーム開始
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
			g.start()
		}
	case gameModeStart:
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
				if g.isFinished() {
					g.finish()
				}
			}
		}
	}

	return true
}

func (g *game) isFinished() bool {
	a := len(g.inputs)
	q := 0
	for _, qs := range g.questions {
		q += len(qs)
	}
	if a >= q {
		return true
	} else {
		return false
	}
}

func (g *game) render() {
	switch g.mode {
	case gameModeStart:
		g.renderMain()
	case gameModeStop:
		g.renderTitle()
	case gameModeFinished:
		g.renderFinished()
	}
}

func (g *game) renderTitle() {
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

func (g *game) renderMain() {
	var current int
	for i, q := range g.questions {
		for j, r := range q {
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault
			inputLen := len(g.inputs)
			if current < inputLen {
				if g.inputs[current] == r {
					fg = termbox.ColorGreen
				} else {
					fg = termbox.ColorYellow
				}
			}
			termbox.SetCell(j, i*2, r, fg, bg)
			if inputLen == current {
				termbox.SetCell(j, i*2+1, '^', termbox.ColorWhite, termbox.ColorDefault)
			} else {
				termbox.SetCell(j, i*2+1, ' ', termbox.ColorWhite, termbox.ColorDefault)
			}
			current++
		}
	}
}

func (g *game) renderFinished() {
	w, h := termbox.Size()

	messages := [...]string{
		"-- GameOver --",
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

func (g *game) Go() {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	render := func(g *game) error {
		g.render()
		return termbox.Flush()
	}

	if err := render(g); err != nil {
		return
	}
	for {
		if !g.update(eventQueue) {
			return
		}
		if err := render(g); err != nil {
			return
		}
	}
}
