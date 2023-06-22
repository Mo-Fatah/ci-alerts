package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Context struct {
	webhook       string
	event         string
	author        string
	commit        string
	commit_url    string
	workflow_name string
	workflow_url  string
}

func NewContext() *Context {
	return &Context{
		webhook:       os.Getenv("webhook"),
		event:         os.Getenv("event"),
		author:        os.Getenv("author"),
		commit:        os.Getenv("commit")[0:6],
		commit_url:    os.Getenv("commit_url"),
		workflow_name: os.Getenv("workflow_name"),
		workflow_url:  os.Getenv("workflow_url"),
	}
}

func main() {
	context := NewContext()
	var mention string
	if context.event == "pr" {
		mention = getAuthorSlackID(context.author)
	} else if context.event == "push" {
		mention = "!channel"
	} else {
		log.Fatal("event type should be specified")
	}

	message := fmt.Sprintf(`{"text": "
	>*CI Failed*\n>*Commit*\n><%s|%s>\n>Workflow Failed: %s\n>Workflow Url: %s\n><%s>"}`,
		context.commit_url, context.commit, context.workflow_name, context.workflow_url, mention)
	body := strings.NewReader(message)
	_, err := http.Post(context.webhook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func getAuthorSlackID(author string) string {
	return "@U04ML7YUSG7"
}
