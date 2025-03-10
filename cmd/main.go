package main

import (
	"fmt"
	"go-transaction/runner"
	"log"
	"net/http"
)

func main() {

	runner, err := runner.New()
	if err != nil {
		log.Fatal(err)
	}

	emailService := NewEmailService(runner)

	mux := http.NewServeMux()
	mux.HandleFunc("/serverone", func(writer http.ResponseWriter, request *http.Request) {
		emailService.SendEmail()
		fmt.Fprintf(writer, "All triggered!\n")
	})

	serverOne := NewServerOne(mux)
	serverTwo := NewServerTwo()

	runner.AddReceivers(serverOne, serverTwo)
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}

}
