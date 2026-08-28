# gazelle_ts integration

This isolated consumer workspace links `gazelle_css` and `gazelle_ts` into one
Gazelle binary. Its generation fixture proves that standard `resolve_regexp`
directives add the mapped CSS contract targets to generated TypeScript
libraries in both root and nested Bazel packages. A checked-in live package
also analyzes the mapped macros and their generated dependency edge.

```sh
bazel test //...
bazel build //live:live
bazel run //:gazelle -- update -mode=diff
```

Keeping the integration in a nested Bazel module prevents the Rust and LLVM
toolchains used by `gazelle_ts` from becoming dependencies of `gazelle_css`.
