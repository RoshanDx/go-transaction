package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ServerOne struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func NewServerOne(mux *http.ServeMux) *ServerOne {
	httpServer := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &ServerOne{
		httpServer: httpServer,
	}
}

func (s ServerOne) Start(ctx context.Context) error {
	fmt.Println(fmt.Sprintf("🚀 Starting server-1"))
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("unable to start server-1: %w", err)
	}

	return nil
}

func (s ServerOne) Stop(ctx context.Context) error {
	fmt.Println(fmt.Sprintf("🛑 Stopping server-1"))
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to shutdown server-1: %w", err)
	}

	return nil
}

type ServerTwo struct {
	httpServer *http.Server
}

func NewServerTwo() *ServerTwo {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/serverone", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, world from server-2\n")
	})

	httpServer := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8081"),
		Handler:      mux1,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &ServerTwo{
		httpServer: httpServer,
	}
}

func (s ServerTwo) Start(ctx context.Context) error {

	fmt.Println(fmt.Sprintf("🚀 Starting server-2"))
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("unable to start server-2: %w", err)
	}

	return nil
}

func (s ServerTwo) Stop(ctx context.Context) error {
	fmt.Println(fmt.Sprintf("🛑 Stopping server-2"))
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to shutdown server-2: %w", err)
	}

	return nil
}
