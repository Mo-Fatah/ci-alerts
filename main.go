package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antchfx/jsonquery"
)

func main() {
	webhook := os.Getenv("webhook")
	github_context := os.Getenv("github_context")
	context, err := NewContext(webhook, github_context)
	if err != nil {
		log.Fatal(err)
	}
	message := buildMessage(context)
	body := strings.NewReader(message)
	_, err = http.Post(context.Webhook, "Content-type: application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

func buildMessage(context *Context) string {
	title := fmt.Sprintf("CI Failed For %s", context.Branch)
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

	failed_action := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Failed Action*\n%s"
		},
	`, context.WorkflowName)
	section += failed_action

	failed_job := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Failed Job*\n%s"
		},
	`, getFailedJob(context.JobsUrl))
	section += failed_job

	commit := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Commit*\n<%s|%s>"
		},
	`, context.CommitUrl, context.Commit[:6])
	section += commit

	action_url := fmt.Sprintf(`
		{
			"type": "mrkdwn",
			"text": "*Workflow Url*\n%s"
		}
	`, context.WorkflowUrl)
	section += action_url

	section += `]}`
	return section
}

func getMention(context *Context) string {
	branch := strings.ToLower(context.Branch)
	if context.TriggeringEvent == "pull_request" {
		return getAuthorSlackID(context.Author)
	} else if branch == "main" || branch == "master" {
		return "<!channel>"
	}
	return ""
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

func getFailedJob(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	jq, err := jsonquery.Parse(strings.NewReader(string(bodyBytes)))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	jobs := jsonquery.FindOne(jq, "jobs").ChildNodes()
	for _, job := range jobs {
		status := job.SelectElement("conclusion").FirstChild.Data
		if status == "failure" {
			return job.SelectElement("name").FirstChild.Data
		}
	}
	return ""
}
