package command

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ListenInterruption(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("signal received", s.String())
		quit <- struct{}{}
	}()
}
