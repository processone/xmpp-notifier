#!/bin/bash
go run /notifier/main.go "${INPUT_SERVER_DOMAIN}" "${INPUT_CORRESPONDANT}" "$INPUT_LOGIN" "${INPUT_PASS}" "${INPUT_SERVER_PORT}" "${INPUT_MESSAGE}" "${INPUT_CORRESPONDENT_IS_ROOM}"

