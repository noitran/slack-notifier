Slack-notifier
====================

<p align="center">
<a href="https://hub.docker.com/r/noitran/slack-notifier"><img src="https://img.shields.io/docker/build/noitran/slack-notifier.svg?style=flat-square" alt="Latest Version"></img></a>
<a href="https://github.com/noitran/opendox/slack-notifier"><img src="https://img.shields.io/github/release/noitran/slack-notifier.svg?style=flat-square" alt="Latest Version"></img></a>
<a href="#"><img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square"></a>
</p>

Dockerized tool to send slack notifications via Slack Incomming Webhooks. Ideal for CI/CD deployments. Docker image available [here](https://hub.docker.com/r/noitran/slack-notifier)

### Example

![alt text](https://github.com/noitran/slack-notifier/blob/master/_demo/preview.png?raw=true)

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

#### Simple example

```console
$ export SLACK_WEBHOOK=https://hooks.slack.com/services/Txxxxxx/Bxxxxxx/xxxxxxxx
$ export SLACK_MESSAGE="Message goes here"
$ docker run -e SLACK_WEBHOOK=$SLACK_WEBHOOK -e SLACK_MESSAGE="hello" noitran/slack-notifier
```

#### Advanced example using gitlab CI

```yaml
variables:
  # ...

stages:
  - version
  - init
  - test
  - build
  - deploy
  - notifyDeploy

# Versionee works with tags attached to project.
# Preferable tag format - v1.0, that is going to be v1.0.x.y
# where v1.0 is tag, x is commits after tag, y - current build number
# ----------------------------------------------------------------------------------------------------------------------
version:
  stage: version
  image: ${CI_REGISTRY}/ci-cd/versionee
  # ...


# Installs composer dependencies
# ----------------------------------------------------------------------------------------------------------------------
init:
  stage: init
  image: ${CI_REGISTRY}/ci-cd/php:latest
  script:
    # ...


# Runs integration and phpunit tests
# ----------------------------------------------------------------------------------------------------------------------
test:
  stage: test
  image: ${CI_REGISTRY}/ci-cd/php:latest
  script:
    - cd src && ./vendor/bin/phpunit
  except:
    - master


# Notify slack that deployment failed
# ----------------------------------------------------------------------------------------------------------------------
notifyDeployFailure:
  stage: notifyDeploy
  image: iocaste/slack-notifier:latest
  when: on_failure
  only:
    - wip
    - master
  before_script:
    - eval export CURRENT_ENV=${CURRENT_ENV}
  script:
    # Set webhook depending on environment
    - export SLACK_WEBHOOK=$([ ${CURRENT_ENV} == 'production' ] && echo ${PROD_SLACK_WEBHOOK_URL} || echo ${DEV_SLACK_WEBHOOK_URL})
    - export SLACK_CHANNEL=$([ ${CURRENT_ENV} == 'production' ] && echo "notify-prod-deploy" || echo "notify-dev-deploy")
    - export SLACK_MESSAGE="Successfully deployed *${CI_PROJECT_NAME}* to *${CURRENT_ENV}* cluster"
    - export BUILD_VERSION=$(cat ${VERSIONEE_FILE})
    - |-
      export SLACK_ATTACHMENT=$(cat <<EOF
      {"color": "danger","pretext": "Commit: ${CI_COMMIT_TITLE}","text": "<${CI_PIPELINE_URL}|Pipeline #${CI_PIPELINE_ID}>","title": ":package: Version tag: ${BUILD_VERSION}","footer": "Project URL: <${CI_PROJECT_URL}|${CI_PROJECT_NAME}>"}
      EOF
      )
    - slack-notifier


# Builds and pushes docker images to registry
# ----------------------------------------------------------------------------------------------------------------------
build:
  stage: build
  only:
    - wip
    - master
  image: docker:latest
  services:
    - docker:dind
  script:
    # ...


# Deploys application containers to development or production cluster
# ----------------------------------------------------------------------------------------------------------------------
deploy:
  stage: deploy
  only:
    - wip
    - master
    # ...

# Notify slack that deployment was successful
# ----------------------------------------------------------------------------------------------------------------------
notifyDeploySuccess:
  stage: notifyDeploy
  image: iocaste/slack-notifier:latest
  only:
    - wip
    - master
  before_script:
    - eval export CURRENT_ENV=${CURRENT_ENV}
  script:
    # Set webhook depending on environment
    - export SLACK_WEBHOOK=$([ ${CURRENT_ENV} == 'production' ] && echo ${PROD_SLACK_WEBHOOK_URL} || echo ${DEV_SLACK_WEBHOOK_URL})
    - export SLACK_CHANNEL=$([ ${CURRENT_ENV} == 'production' ] && echo "notify-prod-deploy" || echo "notify-dev-deploy")
    - export SLACK_MESSAGE="Successfully deployed *${CI_PROJECT_NAME}* to *${CURRENT_ENV}* cluster"
    - export BUILD_VERSION=$(cat ${VERSIONEE_FILE})
    - |-
      export SLACK_ATTACHMENT=$(cat <<EOF
      {
        "color": "#52ac88",
        "pretext": "Commit: ${CI_COMMIT_TITLE}",
        "text": "<${CI_PIPELINE_URL}|Pipeline #${CI_PIPELINE_ID}>",
        "title": ":package: Version tag: ${BUILD_VERSION}",
        "footer": "Project URL: <${CI_PROJECT_URL}|${CI_PROJECT_NAME}>"
      }
      EOF
      )
    - slack-notifier

```
