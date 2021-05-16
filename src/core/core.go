package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)
	for {
		switch <-ch {
		default:
			cancel()
		}
	}
}
