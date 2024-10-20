package main

import (
  // system packages
	"log"
	"os"
	"os/signal"
)

func WaitOnInterrupt() {
	log.Println("Starting interrupt listener...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Waiting on interrupt signal!")
	<-stop

	log.Println("Interrupt signal detected! Ending program...")
}
