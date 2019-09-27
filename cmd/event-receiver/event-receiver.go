package main

import (
	"log"
	"net/http"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/receiver"
)

func main() {
	log.Println("event-receiver")

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

	receive := receiver.Receiver{
		Svc:        svc,
		EventQueue: eventQueue,
		Events:     events,
	}

	mux := http.NewServeMux()
	mux.Handle("/event", http.HandlerFunc(receive.ReceiveEvent))
	log.Fatal(http.ListenAndServe("localhost:8090", mux))
}
