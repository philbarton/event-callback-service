package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetEvents() (map[string][]string, error) {
	eventDir := os.Getenv("EVENT_DIR")
	events := make(map[string][]string)
	files, err := ioutil.ReadDir(eventDir)
	if err != nil {
		return nil, fmt.Errorf("reading files from %s : %v", eventDir, err)
	}
	for _, file := range files {

		if file.IsDir() || file.Name() == "..data" { // skip some crud from kubes config map
			continue
		}

		eventFile := eventDir + "/" + file.Name()
		content, err := ioutil.ReadFile(eventFile)
		if err != nil {
			return nil, fmt.Errorf("reading file content from %s : %v", eventFile, err)
		}
		lines := strings.Split(string(content), "\n")

		events[file.Name()] = lines[:len(lines)-1] // last element is always empty
	}

	log.Println("loaded events")
	return events, nil
}
