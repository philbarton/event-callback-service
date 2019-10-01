package main

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"os"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/sender"
)

func main() {

	callbackQueueCredentials := credentials.NewStaticCredentials(
		os.Getenv("CALLBACK_AWS_ACCESS_KEY_ID"),
		os.Getenv("CALLBACK_AWS_SECRET_ACCESS_KEY"),
		"")

	callbackQueueSvc, err := common.GetSqsService(callbackQueueCredentials)
	if err != nil {
		log.Fatal(err)
	}

	callbackQueue, err := common.GetQueue(callbackQueueSvc, "callback")
	if err != nil {
		log.Fatal(err)
	}

	callbackSender := sender.Sender{
		Svc:           callbackQueueSvc,
		CallbackQueue: callbackQueue,
	}

	callbackSender.SendCallback()
}
