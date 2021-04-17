package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	isShuttingDown bool
	server         http.Server
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server = http.Server{
		Addr: ":8080",
	}

	// Perform application startup.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 5)
		fmt.Fprint(w, "Hello world!")
	})

	// Listen on a different Goroutine so the application doesn't stop here.
	go server.ListenAndServe()

	// Listen for the interrupt signal.
	select {
	case <-ctx.Done():
		// If the interrupt signal is received twice, exit immediately.
		if isShuttingDown {
			os.Exit(1)
		}

		// Restore default behavior on the interrupt signal and notify user of shutdown.
		stop()
		isShuttingDown = true
		fmt.Println("shutting down gracefully, press Ctrl+C again to force")

		// Perform application shutdown...
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println(err)
		}
	}
}
