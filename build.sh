#!/bin/bash

read -n 1 -p "Would you want to export the image file? (Y/n) " ans;
ans=${ans:-Y}

VERSION=$(git describe --tags --abbrev=0)
COMMIT=$(git rev-parse --short HEAD)

DOCKERHUB_REPO="giterhub/slender"

echo ""
echo "Latest: $VERSION / $COMMIT"
echo "Repo: $DOCKERHUB_REPO"

docker build -t "slender-base:$VERSION" -f docker/Dockerfile.base .
docker build -t "$DOCKERHUB_REPO:latest" --build-arg BASE_IMAGE="slender-base:$VERSION" --build-arg VERSION="$VERSION" --build-arg COMMIT="$COMMIT" -f docker/Dockerfile.amd64 .

docker images | grep "$DOCKERHUB_REPO"

case $ans in
  y|Y)
    echo "export the image file..."
    docker image save giterhub/slender:latest | gzip > slender.tar.gz
    echo "export done."
    ;;
esac
