# Buildkite CI

Module tests, basic example, gazelle_ts integration on Bazel 9.0.0/8.6.0; commit hooks.

```mermaid
flowchart TD
  E["PR / main / merge group / manual build"] --> P["Trusted bootstrap: no checkout; verify GitHub PR"]
  P --> F{"Untrusted fork author?"}
  F -->|Yes| A["Pipeline writer approves commit; check stays Pending"]
  F -->|No: same repo or verified Renovate| C["Checkout; pin verified PR merge SHA"]
  A --> C
  C --> U["Upload native Buildkite jobs"]
  U --> T["Repository test matrix"]
  U --> H["Hook checks"]
  T --> S["Required Buildkite status"]
  H --> S
  R["Release / deployment triggers"] --> G["Retained GitHub publishing workflows"]
```

Each Bazel job restores its isolated repository/action cache, runs checks, then saves after success. Hosted queues: `oss` (Linux AMD64), `oss_darwin_arm64` (macOS ARM64). Queue `default` is not used. Tools are checksum-pinned; pnpm follows `packageManager` or npm uses `npm ci`.

## Cache isolation

Bazel jobs use the `oss-ci-v2` registry with server-enforced pipeline/branch scopes.
PRs restore their own branch or main; writes stay in their verified branch.
Main restores only main. GitHub **Prefix third-party fork branch names** must stay
on, so a fork branch named `main` cannot write the real main scope. Cache keys and
PR-controlled YAML are not access controls.

Shared hosted cache volumes are removed. Only repository/action caches persist;
Bazelisk and downloaded tools start fresh per job. The new registry starts empty,
so existing untrusted cache entries are never restored. Initial builds run cold.

Create the registry once in the OSS cluster with `.buildkite/cache-policy.json`.
Do not replace it with the unrestricted default registry. Policy scopes come from
Buildkite's authenticated job claims, not environment variables supplied by jobs.

```mermaid
flowchart LR
  M[Main job] -->|save / restore| MC[Main cache scope]
  P[PR job] -->|save / restore| PC[PR branch cache scope]
  MC -->|read only| P
```

Paste `bootstrap.yml` into pipeline settings; never load approval policy from a PR checkout. Forks require writer approval before checkout or hooks. Only GitHub-verified `renovate[bot]` (ID `29139614`, type `Bot`) bypasses. Each new commit needs approval; API errors and stale PR heads fail closed. Enable third-party forks and publish blocked builds as Pending. Keep workflow tokens and pipeline/cluster secrets disabled.

Build main pushes, PR updates/reopens/base changes, and merge groups; disable tags. Preserve PR exclusions shown in `pipeline.yml`. Manual builds support former workflow-dispatch checks. Replace migrated required GHA checks with `buildkite/gazelle-css` from the Buildkite app only after native checks pass.

Remaining GHA: Release Please and bazel-contrib release/BCR publishing retain their GitHub provenance identities.


## Releases

GHA waits for `release-ready/gazelle-css` on the exact release tag commit. Only non-PR main builds publish that status, after the full BK pipeline finishes. Failed builds stop publication; missing/pending checks time out after 90 minutes. No new secrets or BK token permissions.

The reusable GHA release workflow checks that its checkout still matches the verified SHA, packages the source archive, attests provenance, and publishes BCR. Its duplicate Bazel test run and Bazel caches are disabled. Existing runner exceptions remain on GHA.

```mermaid
flowchart TD
  Main[Main commit] --> BK[BK full test pipeline]
  BK --> Status[Main-only release-ready status]
  Tag[Release tag] --> Gate[GHA checks exact SHA, main ancestry, BK success]
  Status --> Gate
  Gate --> Archive[GHA rechecks SHA and packages source]
  Archive --> Provenance[GHA provenance]
  Provenance --> BCR[GitHub release and BCR]
```
