package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	hook := os.Getenv("webhook")
	//hook, ok := core.GetInput("webhook")
	//if !ok {
	//	log.Fatal("webhook should be provided")
	//}
	//event, ok := core.GetInput("event")
	//if !ok {
	//	log.Fatal("event should be provided")
	//}
	//author, _ := core.GetInput("author")
	event := os.Getenv("event")
	author := ""

	var mention string
	if event == "pr" {
		mention = getAuthorSlackID(author)
	} else if event == "push" {
		mention = "!channel"
	} else {
		log.Fatal("event type should be specified")
	}

	message := fmt.Sprintf("{\"text\":\"Hello, world. <%s> \"}", mention)
	body := strings.NewReader(message)
	_, err := http.Post(hook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func getAuthorSlackID(author string) string {
	return "@U04ML7YUSG7"
}
