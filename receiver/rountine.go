package receiver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ServerOne struct {
	httpServer *http.Server
}

func NewServerOne(httpServer *http.Server) *ServerOne {
	return &ServerOne{httpServer: httpServer}
}

func (r ServerOne) Start(ctx context.Context, serverName string) error {
	fmt.Println(fmt.Sprintf("🚀 Starting %s", serverName))
	if err := r.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.New("unable to start servertwo")
	}
	// Listen for shutdown signal in context
	//<-ctx.Done()
	//fmt.Println("🛑 routine1 received shutdown signal")

	return nil
}

func (r ServerOne) Stop(ctx context.Context, serverName string) error {
	fmt.Println(fmt.Sprintf("🛑 Stopping %s", serverName))
	if err := r.httpServer.Shutdown(ctx); err != nil {
		return errors.New("unable to stop serverone")
	}

	return nil
}

type ServerTwo struct {
	httpServer *http.Server
}

func NewServerTwo(httpServer *http.Server) *ServerTwo {
	return &ServerTwo{httpServer: httpServer}
}

func (r ServerTwo) Start(ctx context.Context, serverName string) error {
	fmt.Println(fmt.Sprintf("🚀 Starting %s", serverName))

	//simulate error
	//return errors.New("❌ servertwo failed to start properly")

	if err := r.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.New("unable to start servertwo")
	}

	return nil
}

func (r ServerTwo) Stop(ctx context.Context, serverName string) error {

	fmt.Println(fmt.Sprintf("🛑 Stopping %s", serverName))

	//simulate error
	//return errors.New("❌ ServerTwo failed to stop properly")

	if err := r.httpServer.Shutdown(ctx); err != nil {
		return errors.New("unable to stop servertwo")
	}

	return nil
}
