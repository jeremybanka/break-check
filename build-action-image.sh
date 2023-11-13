#!/bin/bash

docker buildx build --platform linux/amd64,linux/arm64/v8 -t jeremybanka/break-check:latest --push .
