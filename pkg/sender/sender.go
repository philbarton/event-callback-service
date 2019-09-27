package sender

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"philbarton/event-callback-service/pkg/common"
)

type Sender struct {
	Svc           *sqs.SQS
	CallbackQueue *sqs.GetQueueUrlOutput
}

func (s Sender) SendCallback() {

	log.Println("begin receive and callback")

	for {

		result, err := s.Svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            s.CallbackQueue.QueueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(20),          // 20 seconds
			WaitTimeSeconds:     aws.Int64(common.Poll), // Long Poll
		})

		if err != nil {
			log.Println(fmt.Errorf("reading messages from %s : %v", s.CallbackQueue.String(), err))
			continue
		}

		for _, message := range result.Messages {

			target := message.MessageAttributes["target"]

			if target == nil {
				log.Println(fmt.Errorf("no target message attribute from queue, %s", s.CallbackQueue.String()))
				continue
			}

			targetText := target.StringValue

			queue := s.CallbackQueue
			body := message.Body

			err := callback(*targetText, *body)

			if err != nil {
				log.Println(fmt.Errorf("callback failed for target %s : %v", *target, err))
				continue
			}

			_, err = s.Svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      s.CallbackQueue.QueueUrl,
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				log.Println(fmt.Errorf("failed to delete message %s from %s : %v", *message.MessageId, queue.String(), err))
				continue
			}
		}
	}
}

func callback(targetText string, body string) error {
	log.Printf("Target %s\n", targetText)
	log.Printf("Body %s\n", body)
	/*
		response, err := http.Post(targetText, "text/plain", bytes.NewBufferString(body))

		if err != nil {
			return err
		}

		log.Println(response.Status)
	*/
	return nil
}
