package main

import (
	"log"
	"net/http"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/receiver"
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

	receive := receiver.Receiver{
		Svc:        svc,
		EventQueue: eventQueue,
		Events:     events,
	}

	mux := http.NewServeMux()
	log.Println("begin receive")
	mux.Handle("/event", http.HandlerFunc(receive.ReceiveEvent))
	mux.Handle("/healthz", http.HandlerFunc(receive.Health))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
