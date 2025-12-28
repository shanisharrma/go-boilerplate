package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shanisharrma/go-boilerplate/internal/app"
)

const DefaultContextTimeout = 10

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app, err := app.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	// Start Server
	go func() {
		if err := app.Server.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)
	defer cancel()

	if err = app.Server.Shutdown(ctx); err != nil {
		app.Logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	app.Logger.Info().Msg("server exited properly")

}
