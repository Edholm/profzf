package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()
	if cmd, err := newRootCommand().ExecuteContextC(ctx); err != nil {
		log.Printf("failed to execute %s: %v", cmd.Name(), err)
	}
}
