package receiver

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io/ioutil"
	"log"
	"net/http"
)

type Receiver struct {
	Svc        *sqs.SQS
	EventQueue *sqs.GetQueueUrlOutput
	Events     map[string][]string
}

func (r Receiver) ReceiveEvent(w http.ResponseWriter, req *http.Request) {
	eventType := req.Header["Event-Type"]
	if eventType == nil {
		log.Println(fmt.Errorf("no Event-Type header"))
		w.WriteHeader(400)
		return
	}

	targets := r.Events[eventType[0]]
	if targets == nil {
		log.Println(fmt.Errorf("event not registered %s", eventType[0]))
		w.WriteHeader(400)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(fmt.Errorf("failed to read body : %v", err))
		w.WriteHeader(400)
	} else {
		log.Println(eventType)
		log.Println(string(body))

		attributes := map[string]*sqs.MessageAttributeValue{
			"eventName": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType[0]),
			},
		}
		queue := r.EventQueue

		bodyString := string(body)

		result, err := r.Svc.SendMessage(&sqs.SendMessageInput{
			MessageAttributes: attributes,
			MessageBody:       &bodyString,
			QueueUrl:          queue.QueueUrl,
		})

		if err != nil {
			log.Println(fmt.Errorf("sending message to %s : %v", queue.String(), err))
			w.WriteHeader(500)
			return
		}
		log.Printf("wrote %s\n", *result.MessageId)

		w.WriteHeader(202)
	}
}

func (r Receiver) Health(w http.ResponseWriter, req *http.Request) {

	_, err := r.Svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: r.EventQueue.QueueUrl,
		AttributeNames: []*string{
			aws.String("ApproximateNumberOfMessages"),
		},
	})

	if err != nil {
		log.Println(fmt.Errorf("health check failed : %v", err))
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
	}
}
