package demo

import (
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/storage"
	"os"
)

func PeekQueue() {

	account := os.Getenv("AZURE_ACCOUNT")
	key := os.Getenv("AZURE_KEY")

	cli, _ := storage.NewBasicClient(account, key)

	queueService := cli.GetQueueService()

	eventQueue := queueService.GetQueueReference("event")

	// callbackQueue := queueService.GetQueueReference("callback")

	// NO FREKKING add to queue!!!!!!!!!

	messages, _ := eventQueue.GetMessages(&storage.GetMessagesOptions{})

	for _, message := range messages {
		bytes, _ := base64.URLEncoding.DecodeString(message.Text)
		fmt.Println(string(bytes))
	}

}
