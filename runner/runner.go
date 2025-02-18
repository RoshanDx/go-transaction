package runner

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(serverOne *http.Server, serverTwo *http.Server) error {

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 2)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Starting receivers...")

	wg.Add(2)
	go ServerOneStart(ctx, serverOne, errChan, &wg)
	go ServerTwoStart(ctx, serverTwo, errChan, &wg)

	wg.Add(1)
	go DoWork(ctx, errChan, &wg)

	select {
	case sign := <-signalChan:
		fmt.Println("📡 Signal caught:", sign.String())
		cancel() // send cancellation to all goroutine
	case err := <-errChan:
		fmt.Println("❌ Error caught:", err.Error())
		return err
	}

	currTime := time.Now()

	// create a timeout context for shutting down
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ServerOneStop(shutdownCtx, serverOne, errChan, &wg)
	ServerTwoStop(shutdownCtx, serverTwo, errChan, &wg)

	//Wait for either shutdown completion or timeout
	select {
	case <-shutdownCtx.Done():
		fmt.Printf("⏳ Timeout! Forced shutdown due to hanging runners: %s\n", time.Since(currTime))
	case err := <-errChan:
		fmt.Println("⚠️ Shutdown error:", err)
	default:
		fmt.Println("✅ All runners stopped gracefully")
	}

	wg.Wait()
	close(errChan)
	fmt.Println("✅ Shutdown completed")

	//log remaining shutdown errors
	for err := range errChan {
		fmt.Println("⚠️ Shutdown error:", err)
	}

	return nil
}

func ServerOneStart(ctx context.Context, httpServer *http.Server, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(fmt.Sprintf("🚀 Starting server-1"))
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- errors.New("unable to start server-1")
	}
}

func ServerOneStop(ctx context.Context, httpServer *http.Server, errCh chan<- error, wg *sync.WaitGroup) {
	fmt.Println(fmt.Sprintf("🛑 Stopping server-1"))
	if err := httpServer.Shutdown(ctx); err != nil {
		errCh <- errors.New("unable to shutdown server-1")
	}
}

func ServerTwoStart(ctx context.Context, httpServer *http.Server, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(fmt.Sprintf("🚀 Starting server-2"))
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- errors.New("❌ unable to start server-2")
	}
}

func ServerTwoStop(ctx context.Context, httpServer *http.Server, errCh chan<- error, wg *sync.WaitGroup) {

	// Simulate error or timeout
	//errCh <- errors.New("❌ simulate error from stopping server-2")
	//time.Sleep(13 * time.Second)

	fmt.Println(fmt.Sprintf("🛑 Stopping server-2"))
	if err := httpServer.Shutdown(ctx); err != nil {
		errCh <- errors.New("unable to shutdown server-2")
	}
}

func DoWork(ctx context.Context, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("🚧 Doing work...")
	workerNum := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(3 * time.Second)
			workerNum += 1
			fmt.Println(fmt.Sprintf("👷 Worker-%d completed", workerNum))
		}
		//if workerNum == 3 {
		//	errCh <- errors.New("❌ simulate error from Worker-3")
		//}
	}
}
