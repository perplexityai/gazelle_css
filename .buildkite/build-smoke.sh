#!/usr/bin/env bash
set -euo pipefail
source .buildkite/setup.sh
bazel_command=bazelisk
root=$PWD
for platform in linux-amd64 darwin-arm64; do
  case "$platform" in
    linux-amd64) target_platform=//.buildkite:linux-amd64 ;;
    darwin-arm64) target_platform=//.buildkite:darwin-arm64 ;;
  esac
  echo "--- Build $platform binaries"
  "$bazel_command" build "--platforms=$target_platform" //.buildkite:gazelle //css:css_test
  directory="$root/.buildkite-artifacts/smoke/$USE_BAZEL_VERSION/$platform"
  mkdir -p "$directory"
  for item in 'gazelle //.buildkite:gazelle' 'unit-test //css:css_test'; do
    read -r name target <<< "$item"
    binary=$("$bazel_command" cquery "--platforms=$target_platform" "$target" --output=starlark \
      '--starlark:expr=providers(target)["DefaultInfo"].files_to_run.executable.path')
    cp -L "$binary" "$directory/$name"
  done
  python3 .buildkite/artifacts.py pack "$directory" "$platform"
done
buildkite-agent artifact upload ".buildkite-artifacts/smoke/$USE_BAZEL_VERSION/**/*"
