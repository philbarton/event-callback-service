package main

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"net/http"
	"os"
	"philbarton/event-callback-service/pkg/common"
	"philbarton/event-callback-service/pkg/receiver"
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

	eventQueue, err := common.GetQueue(eventQueueSvc, "event")
	if err != nil {
		log.Fatal(err)
	}

	receive := receiver.Receiver{
		Svc:        eventQueueSvc,
		EventQueue: eventQueue,
		Events:     events,
	}

	mux := http.NewServeMux()
	log.Println("begin receive")
	mux.Handle("/event", http.HandlerFunc(receive.ReceiveEvent))
	mux.Handle("/healthz", http.HandlerFunc(receive.Health))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
