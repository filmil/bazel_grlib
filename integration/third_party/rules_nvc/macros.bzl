load("//build/nvc:rules.bzl", "vhdl_run", "vhdl_elaborate", "vhdl_library")

def vhdl_testbench(name, srcs, deps, entity=None, args=[]):
    vhdl_library_name = "{}_lib".format(name)
    vhdl_library(
        name = vhdl_library_name,
        srcs = srcs,
        deps = deps,
    )
    e = "{}_tb".format(name)
    if entity:
        e = entity
    vhdl_elaborate(
        name = e,
        library = ":{}".format(vhdl_library_name)
    )
    vhdl_run(
        name = name,
        entity = ":{}".format(e),
        args = args,
    )
