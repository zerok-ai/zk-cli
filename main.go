package main

import (
	cmd "github.com/zerok-ai/zk-cli/zkctl/cmd"

	"context"
	"os"
	"os/signal"
)

func main() {
	ctx, cleanup := contextWithSignalInterrupt()
	defer cleanup()

	cmd.ExecuteContext(ctx)
}

func contextWithSignalInterrupt() (context.Context, func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	cleanup := func() {
		signal.Stop(signalChan)
		cancel()
	}

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cleanup
}
