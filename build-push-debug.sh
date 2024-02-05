#!/usr/bin/env bash

set -eo pipefail

THIS_DIR=$(cd "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)

cd "$THIS_DIR"

set -x

export VERSION=2.3.0-dev
export IMAGE_TAG_BASE=ghcr.io/dmrub/piraeus-operator
export IMG=ghcr.io/dmrub/piraeus-operator:v${VERSION}
make docker-push
