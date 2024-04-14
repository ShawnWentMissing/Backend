package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	storage := NewGameStorage()

	startTime := time.Now()
	timeAreas := []TimeArea{
		{1.5, FrontWithinBoundary},
		{3.0, Floor},
		{5.0, Floor},
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueURL := os.Getenv("SQS_URL")
	if queueURL == "" {
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	msgsCh := make(chan TimeArea)

	storage.AddGame("1")

	go pollMessages(floatAreas, startTime, msgsCh)
	go processMessagesSQS(svc, queueURL, msgsCh, storage)

	fmt.Println("Press ctrl + c to exit")
	<-sigCh
	fmt.Println("Exiting")
}
