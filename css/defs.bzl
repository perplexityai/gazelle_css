"""Minimal CSS aggregation rule used by generated BUILD files."""

CssInfo = provider(fields = {"transitive_sources": "All CSS source files."})

def _css_library_impl(ctx):
    transitive = [dep[CssInfo].transitive_sources for dep in ctx.attr.deps]
    files = depset(ctx.files.srcs, transitive = transitive)
    return [
        CssInfo(transitive_sources = files),
        DefaultInfo(files = files),
    ]

css_library = rule(
    implementation = _css_library_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = [".css"]),
        "deps": attr.label_list(providers = [CssInfo]),
    },
)
