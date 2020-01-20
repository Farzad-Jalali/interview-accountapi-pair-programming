#!/bin/bash
echo Packaging dockerfile
source "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/docker-tag.sh"
docker build -t form3tech/$(basename $1):$(git describe) --build-arg APPNAME=$(basename $1) --build-arg TAGS={version:\"$(git describe)\"} -f build/package/$(basename $1)/Dockerfile .
docker tag form3tech/$(basename $1):$(git describe) form3tech/$(basename $1):$TAG