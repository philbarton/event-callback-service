package demo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func ReadWrite() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "philbart",
	}))

	svc := sqs.New(sess)

	eventURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("event"),
	})

	if err != nil {
		fmt.Println("URL Error", err)
		return
	}

	callbackURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("callback"),
	})

	if err != nil {
		fmt.Println("URL Error", err)
		return
	}

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            eventURL.QueueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	for _, message := range result.Messages {
		fmt.Println(message.Body)

		result, err := svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			MessageBody:  message.Body,
			QueueUrl:     callbackURL.QueueUrl,
		})

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		fmt.Println("Success", *result.MessageId)

		resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      eventURL.QueueUrl,
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
			return
		}

		fmt.Println("Message Deleted", resultDelete)
	}
}
