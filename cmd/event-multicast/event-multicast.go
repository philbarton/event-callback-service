package main

import (
	"log"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/multicast"
)

func main() {

	events, err := common.GetEvents()
	if err != nil {
		log.Fatal(err)
	}

	svc, err := common.GetSqsService()
	if err != nil {
		log.Fatal(err)
	}

	eventQueue, err := common.GetQueue(svc, "event")
	if err != nil {
		log.Fatal(err)
	}

	callbackQueue, err := common.GetQueue(svc, "callback")
	if err != nil {
		log.Fatal(err)
	}

	multicaster := multicast.Multicaster{
		Svc:           svc,
		EventQueue:    eventQueue,
		CallbackQueue: callbackQueue,
		Events:        events,
	}

	multicaster.ReceiveAndSend()
}
