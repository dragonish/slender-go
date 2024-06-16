#!/bin/bash

VERSION=$(git describe --tags --abbrev=0)
COMMIT=$(git rev-parse --short HEAD)

DOCKERHUB_REPO="giterhub/slender"

echo "Latest: $VERSION / $COMMIT"
echo "Repo: $DOCKERHUB_REPO"

docker build -t "slender-base:$VERSION" -f docker/Dockerfile.base .
docker build -t "$DOCKERHUB_REPO:latest" --build-arg BASE_IMAGE="slender-base:$VERSION" --build-arg VERSION="$VERSION" --build-arg COMMIT="$COMMIT" -f docker/Dockerfile.amd64 .

docker images | grep "$DOCKERHUB_REPO"
