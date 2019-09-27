#!/usr/bin/env bash

docker run --rm \
    --name callback-sender \
    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
    callback-sender:latest
