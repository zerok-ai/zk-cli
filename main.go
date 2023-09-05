package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"zkctl/cmd"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
)

//go:embed helm-charts/*
//go:embed db-helm-charts/*
//go:embed scripts/*
//go:embed "vizier/vizier.yaml"
var content embed.FS

func main() {
	ui.GlobalWriter.Println("❄ lowering the temperature \n")
	utils.ResetErrorDumpfile()
	ctx, cleanup := contextWithSignalInterrupt()
	defer cleanup()

	var cleanErrorReporter func()
	cleanErrorReporter = ui.InitializeErrorReportor(ctx)
	defer cleanErrorReporter()

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
