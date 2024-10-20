package discord

import (
  // system packages
  "log"
  "os"
  "os/signal"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// start a Discord server connection session
func Start(s *dgo.Session) {
	log.Println("Starting new session...")

	e := s.Open()

	// if error is found
	if e != nil {
		log.Fatalf("fail to start a Discord server connection session [%v]", e)
	}

	log.Println("Successfully started new session!")
}

func Loop(s *dgo.Session) {
	log.Println("Waiting on interrupt signal!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Interrupt signal detected! Ending program...")
}

func Stop(s *dgo.Session) {
	log.Println("Closing session and ending bot...")
	s.Close()
}
