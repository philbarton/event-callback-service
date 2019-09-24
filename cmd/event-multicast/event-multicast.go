package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"philbarton/event-callback-service/pkg/demo"
	"strings"
)

func main() {

	eventDir := os.Getenv("EVENT_DIR")
	events := make(map[string][]string)

	files, err := ioutil.ReadDir(eventDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Println("processing file ", file.Name())
		if file.IsDir() {
			log.Println("Is dir ", file.Name())
			continue
		}

		if file.Name() == "..data" {
			continue
		}

		eventFile := eventDir + "/" + file.Name()
		content, err := ioutil.ReadFile(eventFile)
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(content), "\n")

		events[file.Name()] = lines[:len(lines)-1] // last element is always empty
	}

	fmt.Println(events)

	demo.ReadWrite(events)
}
