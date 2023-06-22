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
	message := buildMessage("CI Failed", context)
	body := strings.NewReader(message)
	_, err := http.Post(context.webhook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func getAuthorSlackID(author string) string {
	return "@U04ML7YUSG7"
}

func buildMessage(title string, context *Context) string {
	var message string
	header := fmt.Sprintf(`{
		"type" : "header",
		"text" : {
			"type": "plain_text",
			"text": "%s"
		}
	},`, title)
	message += header
	message += buildSection(context)
	message = fmt.Sprintf(`{"blocks":[%s]}`, message)
	return message
}

func buildSection(context *Context) string {
	section := `{"type": "section", "fields":[`
	commit := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Commit*\n<%s|%s>"
		},
	`, context.commit_url, context.commit)
	section += commit

	failed_action := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Failed Action*\n%s"
		},
	`, context.workflow_name)
	section += failed_action

	action_url := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Workflow Url*\n%s"
		},
	`, context.workflow_url)
	section += action_url

	var mention string
	if context.event == "pr" {
		mention = getAuthorSlackID(context.author)
	} else if context.event == "push" {
		mention = "!channel"
	} else {
		log.Fatal("event type should be specified")
	}
	mention = fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "<%s>" 
		}
	`, mention)
	section += mention
	section += `]}`
	return section
}
