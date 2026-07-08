#!/bin/sh
set -eu

IMAGE="${IMAGE:-rkeri/goapp-demo-test}"
TAG="${1:-latest}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"

docker build --platform $PLATFORMS -t $IMAGE:$TAG .
