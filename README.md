# gazelle_css

A Gazelle language extension that generates `css_library` targets for CSS
packages. The initial implementation intentionally stays small; this repository
is scaffolded for Bazel testing, examples, automated GitHub releases, and Bazel
Central Registry publication.

## Usage

Add the module to `MODULE.bazel`:

```starlark
bazel_dep(name = "gazelle", version = "0.50.0")
bazel_dep(name = "gazelle_css", version = "0.0.0")
```

Compose the extension into a Gazelle binary in the root `BUILD.bazel`:

```starlark
load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_bin",
    languages = ["@gazelle_css//css"],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_bin",
)
```

Then run `bazel run //:gazelle`. Each package containing ordinary `.css` files
receives one `css_library` named after the package directory. CSS Modules retain
that behavior unless a subtree opts into separate contracts:

```starlark
# gazelle:css_module_enabled true
```

Within an enabled subtree, files ending in `.module.css` are kept out of the
ordinary library and collected in a separate
`css_module_library(name = "css")` rule.

The extension also exposes an abstract `css_module_library` kind for CSS Module
contracts. Its built-in implementation is a source-only fallback; consumers
map it to their contract-producing macro:

```starlark
# gazelle:map_kind css_module_library my_css_module_library //tools:css.bzl
```

Gazelle then adds, updates, and removes the mapped rule and its custom load on
subsequent runs while the plugin continues to reason about `css_module_library`. Use
one direct mapping from the abstract kind to the consumer macro.
Packages with CSS Modules must have one unambiguous aggregate target. Gazelle
fails before rewriting a BUILD file when that target name is already owned or
when multiple module rules exist; configure distinct names or use the standard
whole-file `# gazelle:ignore` escape hatch.
Each directory containing `.module.css` files must also be its own Bazel package;
add a `BUILD.bazel` file instead of collecting nested module sources into a
parent package.

## Directives

| Directive | Default | Purpose |
| --- | --- | --- |
| `css_extension` | `enabled` | Use `disabled` to skip a subtree. |
| `css_library_name` | package basename | Override the generated target name. |
| `css_visibility` | `//visibility:public` | Space-separated visibility labels. |
| `css_module_enabled` | `false` | Set to `true` to generate CSS Modules for a subtree. |
| `css_module_name` | `css` | Name of the aggregate CSS Module target. |

## Development

```sh
bazel test //...
cd examples/basic && bazel run //:gazelle -- update -mode=diff
```

CI tests the module and consumer example on Bazel 8.x and 9.x. Releases use
conventional commits, release-please, `bazel-contrib`'s ruleset release
workflow, and `publish-to-bcr`.

## License

Apache 2.0.
