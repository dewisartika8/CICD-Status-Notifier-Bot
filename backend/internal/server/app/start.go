package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// startServerWithGracefulShutdown function for starting server with a graceful shutdown.
func (s *service) startServerWithGracefulShutdown() {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	// Create context for Telegram bot
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		s.Logger.Info("Shutting down server...")

		// Stop Telegram bot
		if s.TelegramBotManager != nil {
			s.Logger.Info("Stopping Telegram bot...")
			cancel() // Cancel the context to stop bot polling
			s.TelegramBotManager.StopBot()
		}

		// Shutdown HTTP server
		if err := s.HTTPServer.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			s.Logger.Error(fmt.Sprintf("Oops... Server is not shutting down!, err: %s", err.Error()))
		}

		close(idleConnsClosed)
	}()

	// Start Telegram bot in a separate goroutine if available
	if s.TelegramBotManager != nil {
		go func() {
			s.Logger.Info("Starting Telegram bot with polling...")
			s.TelegramBotManager.StartTelegramBot(ctx)
		}()
	}

	// Run server.
	s.Logger.Info(fmt.Sprintf("Starting HTTP server on port %d", s.AppConfig.Server.Port))
	if err := s.HTTPServer.Listen(buildConnectionURL(s.AppConfig.Server.Port)); err != nil {
		s.Logger.Error(fmt.Sprintf("Oops... Server is not running!, err: %s", err.Error()))
	}

	<-idleConnsClosed
}

func buildConnectionURL(port int) string {
	return fmt.Sprintf(":%d", port)
}
