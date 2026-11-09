#!/usr/bin/env bash
set -euo pipefail
source .github/ci/setup.sh
bazel_command=bazelisk
python3 -B .github/ci/artifacts_test.py
root=$PWD
for platform in linux-amd64 darwin-arm64 windows-amd64; do
  case "$platform" in
    linux-amd64) target_platform=//.github/ci:linux-amd64 ;;
    darwin-arm64) target_platform=//.github/ci:darwin-arm64 ;;
    windows-amd64) target_platform=//.github/ci:windows-amd64 ;;
  esac
  echo "--- Build $platform binaries"
  "$bazel_command" build --use_target_platform_for_tests "--platforms=$target_platform" //.github/ci:gazelle //css:css_test
  directory="$root/.ci-artifacts/smoke/$USE_BAZEL_VERSION/$platform"
  mkdir -p "$directory"
  for item in 'gazelle //.github/ci:gazelle' 'unit-test //css:css_test'; do
    read -r name target <<< "$item"
    binary=$("$bazel_command" cquery --use_target_platform_for_tests "--platforms=$target_platform" "$target" --output=starlark \
      '--starlark:expr=providers(target)["DefaultInfo"].files_to_run.executable.path')
    suffix=
    [[ "$platform" != windows-amd64 ]] || suffix=.exe
    cp -L "$binary" "$directory/$name$suffix"
  done
  python3 .github/ci/artifacts.py pack "$directory" "$platform"
done
