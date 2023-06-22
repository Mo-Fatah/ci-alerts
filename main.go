package main

import (
	"bufio"
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
	title := "CI Failed"
	if context.event == "push" {
		title += " On Main"
	}
	message := buildMessage(title, context)
	body := strings.NewReader(message)
	_, err := http.Post(context.webhook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func buildMessage(title string, context *Context) string {
	header := fmt.Sprintf(`{
		"type" : "section",
		"text" : {
			"type": "mrkdwn",
			"text": "*%s* %s"
		}
	},`, title, getMention(context))
	section := buildSection(context)
	message := fmt.Sprintf(`{"blocks" : [ %s ], "attachments":[{ "color": "#a60021", "blocks": [ %s ] }]}`, header, section)
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
		}
	`, context.workflow_url)
	section += action_url

	section += `]}`
	return section
}

func getMention(context *Context) string {
	if context.event == "pr" {
		return getAuthorSlackID(context.author)
	} else if context.event == "push" {
		return "<!channel>"
	}
	panic("event type should be specified")
}

func getAuthorSlackID(author string) string {
	path := os.Getenv("users_path")
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineArr := strings.Split(scanner.Text(), ":")
		if len(lineArr) == 2 {
			if lineArr[0] == author {
				return fmt.Sprintf("<@%s>", lineArr[1])
			}
		}
	}
	os.Exit(0) // exit the program, to prevent any notification from unknown authors' PRs
	return ""
}
