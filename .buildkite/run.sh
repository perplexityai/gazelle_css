#!/usr/bin/env bash
set -euo pipefail
source .buildkite/setup.sh

case "${1:-}" in
  test)
    bazelisk test //...
    ;;
  example)
    cd examples/basic
    bazelisk build //styles
    bazelisk run //:gazelle -- update -mode=diff
    ;;
  integration)
    cd integration/gazelle_ts
    bazelisk test //...
    bazelisk build //live:live
    bazelisk run //:gazelle -- update -mode=diff
    ;;
  *) echo 'Unknown CI suite' >&2; exit 1 ;;
esac
