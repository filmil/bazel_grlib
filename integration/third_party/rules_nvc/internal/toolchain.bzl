load("//internal:providers.bzl", "NVCInfo")

NVC_TOOLCHAIN_TYPE = "@rules_nvc//build/nvc:toolchain_type"
NVC_WRAPPER = Label("@rules_nvc//build/nvc:nvc_wrapper")
NVC_DIRECT_WRAPPER = Label("@rules_nvc//build/nvc:nvc_direct_wrapper")
VHDL_STANDARD_DEFAULT = "2019"

STANDARD_TO_SUFFIX = {
    "1992": "",
    "2008": ".08",
    "2019": ".19",
}


def _nvc_toolchain_impl(ctx):
  toolchain_info = platform_common.ToolchainInfo(
    nvc_info = NVCInfo(
      analyzer = ctx.attr.analyzer,
      artifacts_dir = ctx.attr.artifacts_dir,
      deps = ctx.files.deps,
    ),
  )
  return [toolchain_info]


nvc_toolchain = rule(
  doc = "Defines the NVC toolchain, linking to the NVC analyzer and standard library.",
  implementation = _nvc_toolchain_impl,
  attrs = {
    "analyzer": attr.label(
        executable = True,
        cfg = "host",
        doc = "The NVC executable wrapper script.",
    ),
    "artifacts_dir": attr.label(
        doc = "The directory containing NVC standard libraries and artifacts.",
    ),
    "deps": attr.label_list(
        default = [],
        doc = "Additional toolchain dependencies.",
    ),
  }
)
