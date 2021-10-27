#!/usr/bin/env bash

set -euo pipefail

function usage() {
  echo "Usage: $0 <check|release>"
  exit 1
}

function check() {
  output=${1:-verbose}
  for cmd in git semver; do
    command -v ${cmd} 2>/dev/null 1>&2 || (echo "Error: ${cmd} not found in path" && exit 1)
    [[ "$output" = "silent" ]] || echo "Found ${cmd}"
  done
  [[ "$output" = "silent" ]] || echo "OK"
}

function release {
  git fetch --all

  src_branch=$(git rev-parse --abbrev-ref ${1:-HEAD})
  src_branch_base=${src_branch#origin/*}
  src_branch_topic=${src_branch_base%%/*}
  src_branch_version=${src_branch_base#*/}

  if [[ "$src_branch_topic" != "release" ]]; then
    echo "Unsupported branch type: ${src_branch_topic}"
    exit 1
  fi

  last_tag=$(git describe --abbrev=0 ${src_branch})
  tag=$(semver --range $src_branch_version $last_tag || true)
  if [[ -z "$tag" ]]; then
    tag="$src_branch_version.0"
  else
    tag=$(semver --increment $tag)
  fi

  message=$(
    echo -e "Release $tag\n"
    git log --oneline $last_tag..$src_branch | sed -re 's/^[^ ]+/-/' | (grep -vE "^- (Merge (branch|tag) |Bump)" || [ $? -eq 1 ])
  )

  git tag -a "$tag" -m "$message"
  git push --tags
}

action=${1:-}

case "${action}" in
release)
  check silent
  release
  ;;
check)
  check verbose
  ;;
*)
  echo "Error: invalid action: ${action}"
  usage
  ;;
esac
