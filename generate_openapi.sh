#!/bin/sh
if command -v docker; then
    CONTAINER_ENGINE=docker
elif command -v podman; then
    CONTAINER_ENGINE=podman 
else
    echo "No container engine found. Please install docker or podman."
    exit 1
fi

# generate server stubs
$CONTAINER_ENGINE run --rm -v "${PWD}:/local:Z" docker.io/openapitools/openapi-generator-cli generate \
    -i /local/openapi.yaml \
    -g go-gin-server \
    -o /local/server/openapi