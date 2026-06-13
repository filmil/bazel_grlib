load("//internal:utils.bzl", "get_single_file_from", "get_nvc_deps", "get_nvc_ld_library_path")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
     _NVC_DIRECT_WRAPPER = "NVC_DIRECT_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _impl(ctx):
    nvc_info = ctx.toolchains[_NVC_TOOLCHAIN_TYPE].nvc_info
    nvc_deps = get_nvc_deps(nvc_info)
    analyzer_x = nvc_info.analyzer.files.to_list()[0]
    analyzer = analyzer_x.path
    library_name = ctx.attr.library_name or ctx.attr.name
    container_dir = ctx.actions.declare_directory(
        "{}".format(library_name)
    )

    artifacts = nvc_info.artifacts_dir.files.to_list()
    std_lib_dir = artifacts[1] # hopefully stable...
    
    analyzer_dir = analyzer_x.dirname
    base_dir = analyzer_dir[:-4] if analyzer_dir.endswith("/bin") else analyzer_dir

    targets = ctx.attr.srcs
    srcs = []
    for target in targets:
        files = target.files.to_list()
        srcs += files

    all_libraries = []
    flag_libraries = []
    include_dirs = []
    include_dirs_only = []
    deps_files = []
    seen = []
    deps_hdrs_depset = []
    for target in ctx.attr.deps:
        default_info = target[DefaultInfo]
        deps_files += default_info.files.to_list()
        vhdl_provider = target[VHDLLibraryProvider]
        all_libraries += vhdl_provider.libraries
        deps_hdrs_depset += [vhdl_provider.hdrs]
        for name, path in vhdl_provider.libraries:
            flag_libraries += ["-L", "{path}".format(path=path.path)]
            seen += [name]
        for include_dir in vhdl_provider.includes:
            include_dirs += ["-I", include_dir]
            include_dirs_only += [include_dir]
    for include_dir in ctx.attr.includes:
        include_dirs += ["-I", include_dir]
        include_dirs_only += [include_dir]
    for hdr_file in ctx.files.hdrs:
        hdr_dir = hdr_file.dirname
        include_dirs += ["-I", hdr_dir]
        include_dirs_only += [hdr_dir]

    all_hdrs_files = depset([], transitive=deps_hdrs_depset).to_list()

    ctx.actions.run(
        outputs = [container_dir],
        inputs = srcs + deps_files + ctx.files.hdrs + all_hdrs_files + nvc_deps,
        executable = ctx.executable._direct_wrapper,
        env = {
            "NVC_LD_LIBRARY_PATH": get_nvc_ld_library_path(nvc_info, base_dir, ctx.configuration.default_shell_env),
        },
        arguments = [
          analyzer,
          "--std={}".format(ctx.attr.standard),
          "-L", "{}/nvc".format(std_lib_dir.path),
        ] + flag_libraries + [
          "--work={}:{}/{}".format(
            library_name,
            container_dir.path,
            library_name,
          ),
          "-a",
        ] + include_dirs + [f.path for f in srcs],
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
            includes = include_dirs_only,
            hdrs = depset(ctx.files.hdrs, transitive=deps_hdrs_depset)
        ),
        DefaultInfo(files = depset([container_dir]))
    ]


verilog_library = rule(
    doc = "Compiles Verilog source files into a library using NVC.",
    implementation = _impl,
    attrs = {
        "library_name": attr.string(
            doc = "If the target name is not appropriate as a library name, provide one here",
        ),
        "srcs": attr.label_list(
            allow_files = [".v", ".vh", ".sv", ".svh"],
            doc = "List of Verilog source files.",
        ),
        "hdrs": attr.label_list(
            allow_files = [".v", ".vh", ".sv", ".svh"],
            doc = "List of Verilog header files.",
        ),
        "deps": attr.label_list(
            doc = "List of dependency libraries.",
        ),
        "entities": attr.string_list(
            doc = "List of entities provided by this library.",
        ),
        "standard": attr.string(
            default = _VHDL_STANDARD_DEFAULT,
            doc = "The VHDL standard to use for compilation.",
        ),
        "includes": attr.string_list(
            doc = "list of verilog include directories",
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
