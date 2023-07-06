## CI Alerts
a simple action to customize Slack alerts for CI failures using Slack [Webhooks](https://api.slack.com/messaging/webhooks). The Alerts can be configured mention a slack a channel, if a failure happened on the main/master branch or mention the author if it is a pull request.

## Usage

### Master/Main CI Failure
```yaml
name: Alerts
on:
  workflow_run:
    workflows: [Go]
    types: [completed]

jobs:
  on_push_failure:
    runs-on: ubuntu-latest
    if: github.event.workflow_run.conclusion == 'failure' && github.event.workflow_run.event == 'push' && github.event.workflow_run.head_branch == 'master'
    steps:
      - uses: Mo-Fatah/ci-alerts@v1
        env: 
          webhook: ${{ secrets.SLACK_WEBHOOK }}
          event: push
          commit: ${{ github.sha }}
          commit_url: https://github.com/Mo-Fatah/test-alerts/commit/${{ github.sha }}
          author: ${{ github.actor }}
          workflow_name: ${{ github.event.workflow_run.name }}
          workflow_url: ${{ github.event.workflow_run.html_url}}
```
The above example will trigger the `Alerts` action after the `Go` action is `completed`. If the `Go` action failed and the current branch is `master`, the job will be run. The environment variables below should be provided for the slack message customization. This message will mention the `channel` configured in the webhook to alert the team.


### PR CI Failure

```yaml
name: Alerts
on:
  workflow_run:
    workflows: [Go]
    types: [completed]

jobs:
  on_pull_request_failure:
    runs-on: ubuntu-latest
    if: github.event.workflow_run.conclusion == 'failure' && github.event.workflow_run.event == 'pull_request'
    steps:
      - uses: actions/checkout@v3
      - uses: Mo-Fatah/ci-alerts@v1
        env: 
          webhook: ${{ secrets.SLACK_WEBHOOK }}
          event: pr
          commit: ${{ github.sha }}
          commit_url: https://github.com/Mo-Fatah/test-alerts/commit/${{ github.sha }}
          author: ${{ github.actor }}
          workflow_name: ${{ github.event.workflow_run.name }}
          workflow_url: ${{ github.event.workflow_run.html_url}}
          users_path: ${{github.workspace}}/.github/github-to-slack
```
This will run if the `Go` action failed on a PR. you should provide `users_path`, which is the path to the file that contains the mapping from a github username to a slack [client id](https://api.slack.com/authentication/best-practices#client-id) so the author can be mentioned in the message. If the author doesn't exist in the file, no message will be sent to the slack channel.
The format of the file should be as :
```
github-username:slack-id
```
