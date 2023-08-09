package main

import (
	"fmt"
	"strings"

	"github.com/antchfx/jsonquery"
)

type Context struct {
	Webhook         string
	TriggeringEvent string
	Branch          string
	Author          string
	Commit          string
	CommitUrl       string
	WorkflowName    string
	WorkflowUrl     string
	JobsUrl         string
}

func NewContext(Webhook, github_context string) (*Context, error) {
	jq, err := jsonquery.Parse(strings.NewReader(github_context))
	if err != nil {
		return nil, err
	}

	TriggeringEvent := jsonquery.FindOne(jq, "event/workflow_run/event").FirstChild.Data
	Branch := jsonquery.FindOne(jq, "event/workflow_run/head_branch").FirstChild.Data
	Author := jsonquery.FindOne(jq, "actor").FirstChild.Data
	Commit := jsonquery.FindOne(jq, "sha").FirstChild.Data
	repository := jsonquery.FindOne(jq, "repository").FirstChild.Data
	CommitUrl := fmt.Sprintf("https://github.com/%s/commit/%s", repository, Commit)
	WorkflowName := jsonquery.FindOne(jq, "event/workflow_run/name").FirstChild.Data
	WorkflowUrl := jsonquery.FindOne(jq, "event/workflow_run/html_url").FirstChild.Data
	JobsUrl := jsonquery.FindOne(jq, "event/workflow_run/jobs_url").FirstChild.Data

	return &Context{
		Webhook,
		TriggeringEvent,
		Branch,
		Author,
		Commit,
		CommitUrl,
		WorkflowName,
		WorkflowUrl,
		JobsUrl,
	}, nil
}
