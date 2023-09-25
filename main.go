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

var Version string

func main() {
	printPrePostText := true
	//TODO: Better way to control this?
	for _, arg := range os.Args {
		if arg == "--precise" {
			printPrePostText = false
		}
	}

	if printPrePostText {
		ui.GlobalWriter.Println("❄ lowering the temperature " + "\n")
	}
	utils.ResetErrorDumpfile()
	ctx, cleanup := contextWithSignalInterrupt()
	defer cleanup()

	//var cleanErrorReporter func()
	//cleanErrorReporter = ui.InitializeErrorReportor(ctx)
	//defer cleanErrorReporter()

	cmd.ExecuteContext(ctx)
	if printPrePostText {
		ui.GlobalWriter.Println("\n♨ thawing complete. back to room temperature")
	}
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
