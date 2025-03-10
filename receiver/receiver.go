package receiver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Receiver interface {
	Start(ctx context.Context, serverName string) error
	Stop(ctx context.Context, serverName string) error
}

type ChannelListener interface {
	StartListener(ctx context.Context, errCh chan<- error)
	StopListener(ctx context.Context, errCh chan<- error)
}

type Handler struct {
	receivers []Receiver
	signals   []os.Signal
	channels  []ChannelListener
}

func NewHandler(receivers ...Receiver) (Handler, error) {
	if len(receivers) == 0 {
		return Handler{}, errors.New("receivers must have at least one receiver")
	}

	return Handler{
		receivers: receivers,
		channels:  make([]ChannelListener, 0),
	}, nil
}

func (h Handler) Run() {

	var wg sync.WaitGroup

	// Listen for catchable stop signChan with a buffer
	signChan := make(chan os.Signal, 1)
	errChan := make(chan error, len(h.receivers))

	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Starting receivers...")
	for i, receiver := range h.receivers {
		serverName := "server " + strconv.Itoa(i+1)
		wg.Add(1)
		go func(r Receiver) {
			defer wg.Done()
			if err := r.Start(context.Background(), serverName); err != nil {
				errChan <- err
			}
			//fmt.Println(fmt.Sprintf("wgdone call from starting: %v", serverName))
		}(receiver)
	}

	select {
	case sign := <-signChan:
		fmt.Println("📡 Signal caught:", sign.String())
	case err := <-errChan:
		fmt.Println("❌ Error caught:", err.Error())
	}

	currTime := time.Now()

	//Begin shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Stopping receivers...")
	for i, receiver := range h.receivers {
		serverName := "server " + strconv.Itoa(i+1)
		wg.Add(1)
		go func(r Receiver) {
			defer wg.Done()
			if err := r.Stop(shutdownCtx, serverName); err != nil {
				errChan <- err
			}
			//fmt.Println(fmt.Sprintf("wgdone call from stopping: %v", serverName))
		}(receiver)
	}

	shutdownDone := make(chan struct{})
	go func() {
		wg.Wait()           // Blocks until all receivers call wg.Done()
		close(shutdownDone) // Signal that shutdown completed
	}()

	//Wait for either shutdown completion or timeout
	select {
	case <-shutdownDone:
		fmt.Println("✅ All receivers stopped gracefully")
	case <-shutdownCtx.Done():
		fmt.Printf("⏳ Timeout! Forced shutdown due to hanging receivers: %s\n", time.Since(currTime))
	}

	//wg.Wait()
	close(errChan)

	//log remaining shutdown errors
	for err := range errChan {
		fmt.Println("⚠️ Shutdown error:", err)
	}

	fmt.Println("✅ Shutdown completed")

}
