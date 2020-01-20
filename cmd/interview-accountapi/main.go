package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi"
)

func main() {
	stopServer := make(chan bool, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		<- stop
		fmt.Println("before terminating")
		stopServer <- true
		fmt.Println("after terminating")
	return
	}()

	api.Configure()
	api.Start(stopServer, make(chan bool, 1))
	fmt.Println("End")
}
