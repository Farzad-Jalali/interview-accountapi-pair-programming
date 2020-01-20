package main

import (
	"os"
	"os/signal"

	api "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi"
)

func main() {
	stopServer := make(chan bool, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		stopServer <- <-stop == os.Interrupt
	}()

	api.Configure()
	api.Start(stopServer, make(chan bool))
}
