#!/bin/bash

# Check if version is passed as command-line argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1

REPO="ghcr.io/chaewonkong/msa-link-scraper"

# build image
docker build -t $REPO:$VERSION .

# push image
docker push $REPO:$VERSION

# print success
echo "Docker image pushed to $REPO:$VERSION successfully."