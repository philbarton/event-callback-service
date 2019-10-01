package main

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"os"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/multicast"
)

func main() {

	events, err := common.GetEvents()
	if err != nil {
		log.Fatal(err)
	}

	eventQueueCredentials := credentials.NewStaticCredentials(
		os.Getenv("EVENT_AWS_ACCESS_KEY_ID"),
		os.Getenv("EVENT_AWS_SECRET_ACCESS_KEY"),
		"")

	eventQueueSvc, err := common.GetSqsService(eventQueueCredentials)
	if err != nil {
		log.Fatal(err)
	}

	callbackQueueCredentials := credentials.NewStaticCredentials(
		os.Getenv("CALLBACK_AWS_ACCESS_KEY_ID"),
		os.Getenv("CALLBACK_AWS_SECRET_ACCESS_KEY"),
		"")

	callbackQueueSvc, err := common.GetSqsService(callbackQueueCredentials)
	if err != nil {
		log.Fatal(err)
	}

	eventQueue, err := common.GetQueue(eventQueueSvc, "event")
	if err != nil {
		log.Fatal(err)
	}

	callbackQueue, err := common.GetQueue(callbackQueueSvc, "callback")
	if err != nil {
		log.Fatal(err)
	}

	multicaster := multicast.Multicaster{
		EventSvc:      eventQueueSvc,
		CallbackSvc:   callbackQueueSvc,
		EventQueue:    eventQueue,
		CallbackQueue: callbackQueue,
		Events:        events,
	}

	multicaster.ReceiveAndSend()
}
