#!/usr/bin/env bash

set -euo pipefail

function buildRelease {
  for target in "$@"; do
    set -x
    CGO_ENABLED=0 GOOS=linux go build -v \
        -ldflags='-s -w -extldflags "-static" -X github.com/taisph/go_common/pkg/commonlog.BasePath='${base_path}/' -X github.com/taisph/go_common/pkg/commonversion.Version='${version}' -X github.com/taisph/go_common/pkg/commonversion.Branch='${branch_name}'' \
        -o ${build_output}/${target##*/} ${target}
    set +x
  done
}

function buildDebug {
  for target in "$@"; do
    set -x
    CGO_ENABLED=0 GOOS=linux go build -v \
        -ldflags='-extldflags "-static" -X github.com/taisph/go_common/pkg/commonlog.BasePath='${base_path}/' -X github.com/taisph/go_common/pkg/commonversion.Version='${version}' -X github.com/taisph/go_common/pkg/commonversion.Branch='${branch_name}'' \
        -o ${build_output}/${target##*/} ${target}
    set +x
  done
}

base_path="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
version=${BUILD_VERSION}
branch_name=${BUILD_BRANCH_NAME}
build_variant=${BUILD_VARIANT}
build_output=out
targets="$*"

# Sanity checks.
[[ -n "${targets}" ]] || (echo "$0: no targets specified"; exit 1)

# Build and package.
case "${build_variant}" in
debug)
    buildDebug ${targets}
    ;;
release)
    buildRelease ${targets}
    ;;
*)
    echo "Error: invalid build mode: ${build_variant}"
    exit 1
esac

