def custom_ts_library(name, srcs, deps = [], **kwargs):
    native.filegroup(
        name = name,
        srcs = srcs + deps,
        **kwargs
    )
