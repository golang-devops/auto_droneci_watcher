# auto_droneci_watcher

A very simplistic auto watcher tool for drone-ci (0.5 currently) to watch multiple project drone yamls and perform signing, secrets, etc

## Quick Start

```bash
# this will place the auto_droneci_watcher binary in your $GOPATH/bin directory
go get -u github.com/golang-devops/auto_droneci_watcher
$GOPATH/bin/auto_droneci_watcher -sampleconfig > sample-config.yml
$GOPATH/bin/auto_droneci_watcher -config=sample-config.yml -loglevel=info
```

## Config

### Config layout

```yaml
projects:  
  - repository: my/repo1
    yaml_file: '$GOPATH/src/my/repo1/.drone.yml'
    secrets:
      - plugins/slack SLACK_WEBHOOK=https://hooks.slack.com/services/...
      - plugins/docker,plugins/slack MY_SECRET=MY_SECRET_VALUE
  - repository: my/repo2
    yaml_file: '$GOPATH/src/my/repo2/.drone.yml'
    secrets:
      - plugins/slack SLACK_WEBHOOK=https://hooks.slack.com/services/...
      - plugins/docker,plugins/slack MY_SECRET=MY_SECRET_VALUE
```

### Format of secret line

This is a custom format `image-1,image-2,...,image-n KEY=VALUE`.