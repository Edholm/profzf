package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var errNonZeroExitCode = errors.New("non-zero exit code")

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	if cmd, err := newRootCommand().ExecuteContextC(ctx); err != nil {
		cancel()
		if errors.Is(err, errNonZeroExitCode) {
			os.Exit(1)
		}
		log.Printf("failed to execute %s: %v", cmd.Name(), err)
		os.Exit(1)
	}
	cancel()
}
