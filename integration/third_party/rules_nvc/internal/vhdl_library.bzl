load("//internal:utils.bzl", "get_single_file_from", "get_nvc_deps", "get_nvc_ld_library_path")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
     _NVC_DIRECT_WRAPPER = "NVC_DIRECT_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _vhdl_library(ctx):
    nvc_info = ctx.toolchains[_NVC_TOOLCHAIN_TYPE].nvc_info
    nvc_deps = get_nvc_deps(nvc_info)
    analyzer_x = nvc_info.analyzer.files.to_list()[0]
    analyzer = analyzer_x.path
    library_name = ctx.attr.library_name or ctx.attr.name
    container_dir = ctx.actions.declare_directory(
        "{}".format(library_name)
    )

    artifacts = nvc_info.artifacts_dir.files.to_list()
    std_lib_dir = None
    for artifact in artifacts:
        if artifact.path.endswith("usr/lib/x86_64-linux-gnu/nvc"):
            std_lib_dir = artifact
            break

    if not std_lib_dir:
        # Fallback if specific directory match fails, try to find a known file
        for artifact in artifacts:
             if artifact.path.endswith("usr/lib/x86_64-linux-gnu/nvc/std/STD.STANDARD"):
                 std_lib_dir = artifact.dirname
                 # dirname of std/STD.STANDARD is std, we need the parent 'nvc'
                 # We'll just construct the path directly below
                 break

    # Construct the path to the nvc library directory based on the analyzer path
    # Analyzer is at external/nvc_deb/usr/bin/nvc
    # Library is at external/nvc_deb/usr/lib/x86_64-linux-gnu/nvc
    analyzer_dir = analyzer_x.dirname
    # Assuming path is .../usr/bin
    base_dir = analyzer_dir[:-4] if analyzer_dir.endswith("/bin") else analyzer_dir
    nvc_lib_path = base_dir + "/lib/x86_64-linux-gnu/nvc"

    targets = ctx.attr.srcs
    srcs = []
    for target in targets:
        files = target.files.to_list()
        srcs += files

    all_libraries = []
    flag_libraries = []
    deps_files = []
    seen = []
    vpi_plugins = []
    for target in ctx.attr.deps:
        default_info = target[DefaultInfo]
        deps_files += default_info.files.to_list()
        vhdl_provider = target[VHDLLibraryProvider]
        all_libraries += vhdl_provider.libraries
        if hasattr(vhdl_provider, "vpi_plugins"):
            vpi_plugins.append(vhdl_provider.vpi_plugins)
        for name, path in vhdl_provider.libraries:
            flag_libraries += ["-L", "{path}".format(path=path.path)]
            seen += [name]

    # Add VPI plugins from this target
    vpi_plugins.append(depset(ctx.files.vpi_plugins))


    ctx.actions.run(
        outputs = [container_dir],
        inputs = srcs + deps_files + nvc_deps,
        executable = ctx.executable._direct_wrapper,
        env = {
            "NVC_LD_LIBRARY_PATH": get_nvc_ld_library_path(nvc_info, base_dir, ctx.configuration.default_shell_env),
        },
        arguments = [
          analyzer,
          "--std={}".format(ctx.attr.standard),
          "-L", nvc_lib_path,
        ] + flag_libraries + [
          "--work={}:{}/{}".format(
            library_name,
            container_dir.path,
            library_name,
          ),
          "-a",
        ] + [f.path for f in srcs],
        tools = [analyzer_x] + artifacts,
        # Only seems to work from bazel 6.0.0 on.
        #toolchain = _NVC_TOOLCHAIN_TYPE,
        progress_message = "VHDL analyze: {}".format(library_name),
    )
    return [
        VHDLLibraryProvider(
            libraries = [(library_name, container_dir)] + all_libraries,
            entities = ctx.attr.entities,
            library_name = library_name,
            library_dir = container_dir,
            vpi_plugins = depset(transitive = vpi_plugins),
        ),
        DefaultInfo(files = depset([container_dir]))
    ]

vhdl_library = rule(
    doc = "Compiles VHDL source files into a library using NVC.",
    implementation = _vhdl_library,
    attrs = {
        "library_name": attr.string(
            doc = "If the target name is not appropriate as a library name, provide one here",
        ),
        "srcs": attr.label_list(
            allow_files = [".vhdl", ".vhd"],
            doc = "A list of VHDL source files.",
        ),
        "deps": attr.label_list(
            doc = "A list of other `vhdl_library` targets that this library depends on.",
        ),
        "entities": attr.string_list(
            doc = "A list of VHDL entities provided by this library.",
        ),
        "standard": attr.string(
            default = _VHDL_STANDARD_DEFAULT,
            doc = "The VHDL standard to use for compilation (e.g., '2008', '2019'). Defaults to '2019'.",
        ),
        "vpi_plugins": attr.label_list(
            allow_files = True,
            doc = "List of VPI plugins required for simulation.",
        ),
        "_direct_wrapper": attr.label(
            default = _NVC_DIRECT_WRAPPER,
            executable = True,
            cfg = "host",
        ),
    },
    toolchains = [
        _NVC_TOOLCHAIN_TYPE,
    ],
)

