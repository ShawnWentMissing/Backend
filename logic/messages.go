package backend

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Message struct {
	ID      string `json:"id"`
	Bounced bool   `json:"bounced"`
	Area    Area   `json:"area"`
}

type TimeArea struct {
	Time float64
	Area Area
}

func pollMessagesSQS(svc *sqs.SQS, queueURL string, msgsCh chan<- *sqs.Message) {
	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(240),
			WaitTimeSeconds:     aws.Int64(1),
		})
		if err != nil {
			fmt.Println("Error receiving messages:", err)
			continue
		}

		for _, msg := range result.Messages {
			msgsCh <- msg
		}

		time.Sleep(1 * time.Second)
	}
}

func processMessagesSQS(svc *sqs.SQS, queueURL string, msgsCh <-chan *sqs.Message, storage *GameStorage) {
	for msg := range msgsCh {
		var decodedMsg Message
		err := json.Unmarshal([]byte(*msg.Body), &decodedMsg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})

		if decodedMsg.Bounced {
			endRally, handout, ok := storage.BallBounce(decodedMsg.ID, decodedMsg.Area)
			if !ok {
				fmt.Println("Error updating game")
				return
			}

			game, ok := storage.GetGame(decodedMsg.ID)
			if !ok {
				fmt.Println("Error updating game")
				return
			}

			if endRally {
				announceMessage(decodedMsg.ID, game.Player1Score, game.Player2Score, handout)
			}
		}
	}
}

func processMessages(msgsCh <-chan Message, storage *GameStorage) {
	for msg := range msgsCh {
		if msg.Bounced {
			endRally, handout, ok := storage.BallBounce(msg.ID, msg.Area)
			if !ok {
				fmt.Println("Error updating game")
				return
			}

			if endRally {
				game, ok := storage.GetGame("1")
				if !ok {
					fmt.Println("Error updating game")
					return
				}

				if endRally {
					announceMessage("1", game.Player1Score, game.Player2Score, handout)
				}
			}
		}
	}
}

func pollMessages(timeAreas []TimeArea, msgsCh chan<- Message) {
	for _, ta := range timeAreas {
		duration := time.Duration(ta.Time * float64(time.Second))
		time.Sleep(duration)
		msgsCh <- Message{"1", true, ta.Area}
	}
}

func announceMessage(id string, player1score, player2score int, handout bool) {
	WebSocketHandler.NotifyClients(message)
}
