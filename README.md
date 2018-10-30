# Slack-notifier

Dockerized tool to send slack notifications via Slack Incomming Webhooks. Ideal for CI/CD deployments.

## Slack Configuration

1. Sign in to your Slack team and start a new [Incoming WebHooks configuration](https://my.slack.com/services/new/incoming-webhook/).
2. Select the Slack channel where notifications will be sent to by default. Click the **Add Incoming WebHooks integration** button to add the configuration
3. Copy the **Webhook URL**


## Build


### Local build

```console
$ make build
```

## Usage

```console
$ export SLACK_WEBHOOK=https://hooks.slack.com/services/Txxxxxx/Bxxxxxx/xxxxxxxx
$ export SLACK_MESSAGE="Message goes here"
$ docker run -e SLACK_WEBHOOK=$SLACK_WEBHOOK -e SLACK_MESSAGE="hello" iocaste/slack-notifier
```


## Available environment variables

| Variable  | Required | Example |
| --------- | -------- | ------- |
| SLACK_WEBHOOK | **Yes** | https://hooks.slack.com/services/Txxxxxx/Bxxxxxx/xxxxxxxx |
| SLACK_USERNAME | No  | DeploymentBot |
| SLACK_CHANNEL | No  | notify-prod-deploy |
| SLACK_ICON | No | http://example.com/icon.png |
| SLACK_MESSAGE | No  | Message body |
| SLACK_ATTACHMENT | No  | Attachment body as json string |


## Example deployment

**gitlab-ci.yml**

```yaml
stages:
  - notifyDeployProd
  
variables:
  SLACK_WEBHOOK: https://hooks.slack.com/services/Txxxxxx/Bxxxxxx/xxxxxxxx
  SLACK_CHANNEL: notify-prod-deploy

notifyDeployProd:
  stage: notifyDeployProd
  image: iocaste/slack-notifier:latest
  only:
  - master
  script:
    - 'SLACK_MESSAGE="Message" slack-notifier'

```
