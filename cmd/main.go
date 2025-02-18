package main

import (
	"fmt"
	"go-transaction/runner"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/serverone", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, world from server-1\n")
	})

	httpServer1 := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      mux1,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/servertwo", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, world from server-2\n")
	})

	httpServer2 := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8081"),
		Handler:      mux2,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := runner.Run(httpServer1, httpServer2); err != nil {
		log.Fatal(err)
	}
}
