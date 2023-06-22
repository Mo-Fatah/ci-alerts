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
	event := os.Getenv("event")
	author := os.Getenv("author")
	commit := os.Getenv("commit")[0:6]
	commit_url := os.Getenv("commit_url")
	workflow_name := os.Getenv("workflow_name")
	workflow_url := os.Getenv("workflow_url")

	var mention string
	if event == "pr" {
		mention = getAuthorSlackID(author)
	} else if event == "push" {
		mention = "!channel"
	} else {
		log.Fatal("event type should be specified")
	}

	message := fmt.Sprintf(`{"text": "
	> *CI Failed*  
	> *Commit* 
	> <%s|%s>
	> Workflow Failed: %s
	> Workflow Url: %s
	> <%s>"}`,
		mention, commit_url, commit, workflow_name, workflow_url)
	body := strings.NewReader(message)
	_, err := http.Post(hook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func getAuthorSlackID(author string) string {
	return "@U04ML7YUSG7"
}
