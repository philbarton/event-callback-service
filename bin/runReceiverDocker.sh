#!/usr/bin/env bash

docker run --rm \
    --name event-receiver \
    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
    -e EVENT_DIR=/events \
    -v /Users/philbarton/tmp/events:/events \
    -p8090:8090 \
    event-receiver:latest
