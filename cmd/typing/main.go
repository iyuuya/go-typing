package main

import (
	"context"
	"time"

	"github.com/iyuuya/go-typing/internal/typing"
)

func main() {
	typing.Setup()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	typing.MainLoop(ctx)

}
