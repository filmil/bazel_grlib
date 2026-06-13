load("//internal:utils.bzl", "get_single_file_from")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _extract_file(ctx):
    filter = []
    for file in ctx.attr.src.files.to_list():
        if ctx.attr.filter in file.path:
            filter += [file]
    return [
        DefaultInfo(files=depset(filter),
        runfiles=ctx.runfiles(files=filter)),
    ]


extract_file = rule(
    doc = "Extracts specific files from a source target based on a string filter.",
    implementation = _extract_file,
    attrs = {
        "filter": attr.string(
            doc = "A string used to filter the files in the source target.",
        ),
        "src" : attr.label(
            doc = "The source target to extract files from.",
        ),
    },
)
