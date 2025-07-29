package app

import (
	"fmt"
	"os"
	"os/signal"
)

// startServerWithGracefulShutdown function for starting server with a graceful shutdown.
func (s *service) startServerWithGracefulShutdown() {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := s.HTTPServer.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			s.Logger.Error(fmt.Sprintf("Oops... Server is not shutting down!, err: %s", err.Error()))
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := s.HTTPServer.Listen(buildConnectionURL(s.AppConfig.Server.Port)); err != nil {
		s.Logger.Error(fmt.Sprintf("Oops... Server is not running!, err: %s", err.Error()))
	}

	<-idleConnsClosed
}

func buildConnectionURL(port int) string {
	return fmt.Sprintf(":%d", port)
}
