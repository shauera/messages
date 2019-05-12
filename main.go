package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	rest "github.com/shauera/messages/rest"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func init() {
	InitConfig()
	InitLogging()
}

func main() {
	log.Println("Starting the messages service...")

	cancellableContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	go rest.StartHTTPServer(cancellableContext)

	// -- Wait for a SIGINT. Run cleanup when signal is received ---
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Println("Received an interrupt, stopping service...")
		// Do some cleanup...

		// Stop all running routines
		cancel()

		// Wait for routines to stop
		shutdownGraceDuration := config.GetDuration("service.shutdownGraceDuration")
		log.WithField("gracePeriod", shutdownGraceDuration).Info("Waiting for routines to stop")
		time.Sleep(shutdownGraceDuration * time.Second)

		// Finished cleanup...
		close(cleanupDone)
	}()
	<-cleanupDone
	// -------------------------------------------------------------
}
