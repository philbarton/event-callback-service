package demo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"strconv"
	"strings"
)

func ReadWrite() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})

	if err != nil {
		log.Println("Failed to establish session", err)
		return
	}

	svc := sqs.New(sess)

	eventURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("event"),
	})

	if err != nil {
		log.Println("Failed to get event queue", err)
		return
	}

	callbackURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("callback"),
	})

	if err != nil {
		log.Println("Failed to get callback queue", err)
		return
	}

	for {

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
			WaitTimeSeconds:     aws.Int64(10), // Long poll
		})

		if err != nil {
			log.Println("Error receiving message ", err)
			continue
		}

		for _, message := range result.Messages {
			log.Println(message.Body)

			for i := range [2]int{} { // simulate multiple consumers

				log.Println(i)
				var sb strings.Builder
				sb.WriteString(*message.Body)
				sb.WriteString(" ")
				sb.WriteString(strconv.Itoa(i))

				newMessage := sb.String()

				result, err := svc.SendMessage(&sqs.SendMessageInput{
					DelaySeconds: aws.Int64(10),
					MessageBody:  &newMessage,
					QueueUrl:     callbackURL.QueueUrl,
				})

				if err != nil {
					log.Println("Error sending", err)
					continue
				}

				log.Println("Success", *result.MessageId)

				resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      eventURL.QueueUrl,
					ReceiptHandle: message.ReceiptHandle,
				})

				if err != nil {
					log.Println("Delete Error....would be a duplicate!", err)
					continue
				}

				log.Println("Message Deleted", resultDelete)
			}
		}
	}
}
