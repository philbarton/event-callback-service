package main

import (
	"log"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/sender"
)

func main() {

	svc, err := common.GetSqsService()
	if err != nil {
		log.Fatal(err)
	}

	callbackQueue, err := common.GetQueue(svc, "callback")
	if err != nil {
		log.Fatal(err)
	}

	callbackSender := sender.Sender{
		Svc:           svc,
		CallbackQueue: callbackQueue,
	}

	callbackSender.SendCallback()
}
