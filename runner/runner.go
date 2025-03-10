package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type JobType string

// JobRunner defines the interface for background job execution.
type JobRunner interface {
	RunJob(ctx context.Context, jobName JobType, fn func())
}

// Receiver defines the interface for running receivers in Run
type Receiver interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Runner struct {
	wg        sync.WaitGroup
	receivers []Receiver
}

func New(options ...func(runner *Runner) error) (*Runner, error) {
	runner := &Runner{}
	for _, option := range options {
		if err := option(runner); err != nil {
			return nil, err
		}
	}
	return runner, nil
}

func WithReceiver(receivers ...Receiver) func(*Runner) error {
	return func(runner *Runner) error {
		if receivers == nil {
			return errors.New("receivers is empty")
		}
		runner.receivers = receivers
		return nil
	}
}

func NewWithRunners(receivers ...Receiver) (*Runner, error) {
	if len(receivers) == 0 {
		return nil, errors.New("receivers is empty")
	}
	return &Runner{
		receivers: receivers,
	}, nil
}

func (r *Runner) AddReceivers(receivers ...Receiver) {
	r.receivers = append(r.receivers, receivers...)
}

func (r *Runner) Run() error {

	if len(r.receivers) == 0 {
		return errors.New("receivers is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, len(r.receivers))

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Starting receivers...")

	for _, receiver := range r.receivers {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			if err := receiver.Start(ctx); err != nil {
				errChan <- err
			}
		}()
	}

	select {
	case sign := <-signalChan:
		fmt.Println("📡 Signal caught:", sign.String())
		cancel()
	case err := <-errChan:
		fmt.Println("❌ Error caught:", err.Error())
		return err
	}

	currTime := time.Now()

	// create a timeout context for shutting down
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, receiver := range r.receivers {
		if err := receiver.Stop(shutdownCtx); err != nil {
			errChan <- err
		}
	}

	//Wait for either shutdown completion or timeout
	select {
	case <-shutdownCtx.Done():
		fmt.Printf("⏳ Timeout! Forced shutdown due to hanging runners: %s\n", time.Since(currTime))
	case err := <-errChan:
		fmt.Println("⚠️ Shutdown error:", err)
	default:
		fmt.Println("✅ All runners stopped gracefully")
	}

	fmt.Println("👷 Completing remaining job runner")
	r.wg.Wait()
	fmt.Println("✅ All job runner stopped gracefully")

	close(errChan)

	fmt.Println("🛑 Shutdown completed")

	//log remaining errors
	for err := range errChan {
		fmt.Println("⚠️ Runner error:", err)
	}

	return nil
}

func (r *Runner) RunJob(ctx context.Context, jobName JobType, fn func()) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Sprintf("background job fail to execute: %v", err))
			}
		}()

		// run job
		fmt.Printf("🚧 Background job started: %s\n", jobName)
		fn()
		fmt.Printf("🚧 Background job completed: %s\n", jobName)

	}()
}
