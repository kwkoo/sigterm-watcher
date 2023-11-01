package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	log.Printf("%s started - waiting for signal...", os.Args[0])
	<-ctx.Done()
	log.Print("signal received, shutting down")
	stop()
}
