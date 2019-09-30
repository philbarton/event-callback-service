package multicast

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"philbarton/event-callback-service/pkg/common"
)

type Multicaster struct {
	Svc           *sqs.SQS
	EventQueue    *sqs.GetQueueUrlOutput
	CallbackQueue *sqs.GetQueueUrlOutput
	Events        map[string][]string
}

func (m Multicaster) ReceiveAndSend() {

	log.Println("begin receive and send")

	for {

		result, err := m.Svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            m.EventQueue.QueueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(20),          // 20 seconds
			WaitTimeSeconds:     aws.Int64(common.Poll), // Long poll
		})

		if err != nil {
			log.Println(fmt.Errorf("reading messages from %s : %v", m.EventQueue.String(), err))
			continue
		}

		for _, message := range result.Messages {

			log.Printf("Read %s\n", *message.MessageId)

			eventName := message.MessageAttributes["eventName"]

			if eventName == nil {
				log.Println(fmt.Errorf("no eventName message attribute from queue, %s", m.EventQueue.String()))
				continue
			}

			event := eventName.StringValue
			targets := m.Events[*event]

			for _, target := range targets {

				attributes := map[string]*sqs.MessageAttributeValue{
					"target": {
						DataType:    aws.String("String"),
						StringValue: aws.String(target),
					},
				}
				queue := m.CallbackQueue
				body := message.Body

				result, err := m.Svc.SendMessage(&sqs.SendMessageInput{
					MessageAttributes: attributes,
					MessageBody:       body,
					QueueUrl:          queue.QueueUrl,
				})

				if err != nil {
					log.Println(fmt.Errorf("sending message to %s : %v", queue.String(), err))
					continue
				}
				log.Printf("wrote %s\n", *result.MessageId)

				_, err = m.Svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      m.EventQueue.QueueUrl,
					ReceiptHandle: message.ReceiptHandle,
				})

				if err != nil {
					log.Println(fmt.Errorf("failed to delete message %s from %s : %v", *message.MessageId, queue.String(), err))
					continue
				}
			}
		}
	}
}
