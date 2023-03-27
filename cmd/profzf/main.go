package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cf()
	if cmd, err := newRootCommand().ExecuteContextC(ctx); err != nil {
		log.Printf("failed to execute %s: %v", cmd.Name(), err)
	}
}
