# Container image that runs your code
FROM golang:1.19

COPY entrypoint.sh /notifier/entrypoint.sh
COPY go.mod /notifier/go.mod
COPY go.sum /notifier/go.sum
COPY main.go /notifier/main.go

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/notifier/entrypoint.sh"]
