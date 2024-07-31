#!/bin/bash
cd "$(dirname "$0")"
go get -d gosrc.io/xmpp
go install gosrc.io/xmpp
go run /notifier/main.go "${INPUT_SERVER_HOST}" "${INPUT_RECIPIENT}" "${INPUT_JID}" "${INPUT_PASSWORD}" "${INPUT_SERVER_PORT}" "${INPUT_MESSAGE}" "${INPUT_RECIPIENT_IS_ROOM}" "${INPUT_BOT_ALIAS}"
