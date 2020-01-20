#!/bin/bash

source "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/docker-tag.sh"

function publish() {
    echo "Publising to docker hub"

    echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin


    docker push form3tech/$NAME:$TAG
    docker push form3tech/$NAME:$VERSION
}

NAME=$(basename $1)
VERSION=$(git describe)

publish