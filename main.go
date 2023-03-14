package main

import (
	"zkctl/cmd"

	"context"
	"os"
	"os/signal"
	"zkctl/cmd/pkg/ui"
)

func main() {
	ui.GlobalWriter.Println("❄ lowering the temperature \n")
	ctx, cleanup := contextWithSignalInterrupt()
	defer cleanup()

	cmd.ExecuteContext(ctx)

	ui.GlobalWriter.Println("\n♨ thawing complete. back to room temperature")
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
