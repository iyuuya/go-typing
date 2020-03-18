package typing

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var words []string

func Setup() {
	words = []string{
		"hello",
		"world",
		"golang",
		"root",
		"insert",
		"update",
		"vim",
		"emacs",
		"word",
		"string",
	}
}

func input(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- string(s.Bytes())
		}
		close(ch)
	}()
	return ch
}

func getRandomWord() string {
	t := time.Now().UnixNano()
	rand.Seed(t)
	i := rand.Intn(len(words))
	return words[i]
}

type result struct {
	ok int
	ng int
}

func (r *result) record(res bool) {
	if res {
		r.ok++
	} else {
		r.ng++
	}
}

func MainLoop(ctx context.Context) {
	res := result{}

	ch := input(os.Stdin)

	go func() {
		for {
			word := getRandomWord()
			fmt.Printf("==> %s\n", word)
			fmt.Print(">")
			inputWord := <-ch
			res.record(inputWord == word)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("OK: %d, NG: %d\n", res.ok, res.ng)
	}
}
