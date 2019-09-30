package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

const Poll = 2

func GetSqsService() (*sqs.SQS, error) {
	region := "eu-west-2"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("create session, region is %s : %v", region, err)
	}
	svc := sqs.New(sess)

	log.Println("got sqs service")
	return svc, err
}

func GetQueue(svc *sqs.SQS, queueName string) (*sqs.GetQueueUrlOutput, error) {
	eventQueue, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get %s queue: %v", queueName, err)
	}
	log.Printf("got %s queue\n", queueName)
	return eventQueue, nil
}
