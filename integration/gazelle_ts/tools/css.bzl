def custom_css_module_library(name, srcs, **kwargs):
    native.filegroup(
        name = name,
        srcs = srcs,
        **kwargs
    )
    native.filegroup(
        name = name + ".web",
        srcs = srcs,
    )
