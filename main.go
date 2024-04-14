package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
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

	msgsCh := make(chan *sqs.Message)

	storage := NewGameStorage()

	storage.AddGame("1")

	go pollMessagesSQS(svc, queueURL, msgsCh)
	go processMessagesSQS(svc, queueURL, msgsCh, storage)

	fmt.Println("Press ctrl + c to exit")
	<-sigCh
	fmt.Println("Exiting")
}
