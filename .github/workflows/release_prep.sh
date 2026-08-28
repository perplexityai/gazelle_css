#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

TAG="${1:-${GITHUB_REF_NAME:-}}"
if [[ -z "${TAG}" ]]; then
    echo "release_prep.sh: tag is required" >&2
    exit 1
fi

VERSION="${TAG#v}"
MODULE=gazelle_css
PREFIX="${MODULE}-${VERSION}"
ARCHIVE="${MODULE}-${TAG}.tar.gz"

git archive --format=tar --prefix="${PREFIX}/" "${TAG}" | gzip -9 > "${ARCHIVE}"

SHA256_HEX=$(shasum -a 256 "${ARCHIVE}" | awk '{print $1}')
SHA256_B64=$(printf '%s' "${SHA256_HEX}" | xxd -r -p | base64)
INTEGRITY="sha256-${SHA256_B64}"

cat <<EOF
## Using Bzlmod

\`\`\`starlark
bazel_dep(name = "${MODULE}", version = "${VERSION}")
\`\`\`

## Using a non-registry override

\`\`\`starlark
bazel_dep(name = "${MODULE}", version = "${VERSION}")
archive_override(
    module_name = "${MODULE}",
    integrity = "${INTEGRITY}",
    strip_prefix = "${PREFIX}",
    urls = ["https://github.com/${GITHUB_REPOSITORY:-perplexityai/${MODULE}}/releases/download/${TAG}/${ARCHIVE}"],
)
\`\`\`
EOF
