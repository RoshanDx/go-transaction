package main

import (
	"context"
	"fmt"
	"go-transaction/runner"
	"time"
)

const EmailJobType = "EmailJobType"

type EmailService struct {
	jobRunner runner.JobRunner
}

func NewEmailService(jobRunner runner.JobRunner) EmailService {
	return EmailService{
		jobRunner: jobRunner,
	}
}

func (s EmailService) SendEmail() {

	emailList := []string{
		"genji@test.com",
		"kaido@test.com",
		"zoro@test.com",
		"robin@test.com",
		"sensei@test.com",
		"enzo@test.com",
		"chopper@test.com",
		"anderson@test.com",
		"richter@test.com",
		"trevor@test.com",
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.jobRunner.RunJob(ctx, EmailJobType, func() {
		for _, email := range emailList {
			select {
			case <-ctx.Done():
				fmt.Println("Work canceled due to parent context.")
				return
			default:
				time.Sleep(2 * time.Second)
				//if i == 3 {
				//	fmt.Printf("Send email fail for %s\n", email)
				//	//cancel()
				//	return
				//}
				fmt.Println("📧 Sending email to " + email)
			}
		}
		cancel()
	})
}
