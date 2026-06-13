load("//internal:utils.bzl", "get_single_file_from", "get_nvc_deps", "get_nvc_ld_library_path")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _vhdl_elaborate(ctx):
    nvc_info = ctx.toolchains[_NVC_TOOLCHAIN_TYPE].nvc_info
    nvc_deps = get_nvc_deps(nvc_info)
    analyzer_x = nvc_info.analyzer.files.to_list()[0]
    analyzer = analyzer_x.path

    vhdl_provider = ctx.attr.library[VHDLLibraryProvider]
    library_name = vhdl_provider.library_name
    out_dir = ctx.actions.declare_directory("{}_out".format(library_name))

    artifacts = nvc_info.artifacts_dir.files.to_list()
    std_lib_dir = None
    for artifact in artifacts:
        if artifact.path.endswith("usr/lib/x86_64-linux-gnu/nvc"):
            std_lib_dir = artifact
            break
            
    if not std_lib_dir:
        for artifact in artifacts:
             if artifact.path.endswith("usr/lib/x86_64-linux-gnu/nvc/std/STD.STANDARD"):
                 std_lib_dir = artifact.dirname
                 break
                 
    analyzer_dir = analyzer_x.dirname
    base_dir = analyzer_dir[:-4] if analyzer_dir.endswith("/bin") else analyzer_dir
    nvc_lib_path = base_dir + "/lib/x86_64-linux-gnu/nvc"


    all_libraries = []
    flag_libraries = []
    deps_paths = []
    seen = []

    vhdl_provider = ctx.attr.library[VHDLLibraryProvider]
    all_libraries += vhdl_provider.libraries
    for name, path in vhdl_provider.libraries:
        if name != vhdl_provider.library_name and name not in seen:
            flag_libraries += ["-L", path.path]
            deps_paths += [path]
            seen += [name]

    runfiles = ctx.runfiles(files=deps_paths + [out_dir])

    runfiles.merge_all([ctx.attr.library[DefaultInfo].default_runfiles])

    work_library_file = get_single_file_from(ctx.attr.library)
    vpi_plugins = []
    vpi_flags = []
    if hasattr(vhdl_provider, "vpi_plugins") and vhdl_provider.vpi_plugins:
        vpi_plugins = vhdl_provider.vpi_plugins.to_list()
        for p in vpi_plugins:
            vpi_flags.append("--load=" + p.path)

    ctx.actions.run(
        outputs = [out_dir],
        inputs = depset(direct = [vhdl_provider.library_dir] + ([std_lib_dir] if hasattr(std_lib_dir, "path") else []) + artifacts + nvc_deps + deps_paths + vpi_plugins).to_list(),
        executable = ctx.executable._script.path,
        env = {
            "NVC_LD_LIBRARY_PATH": get_nvc_ld_library_path(nvc_info, base_dir, ctx.configuration.default_shell_env),
        },
        arguments = [
            "--vhdl-standard={}".format(ctx.attr.standard),
            "--nvc-binary-path={}".format(analyzer),
            "--library-name={}".format(library_name),
            "--library-paths={}".format(" ".join(flag_libraries + ["-L", nvc_lib_path])),
            "--stdlib-dir={}".format(nvc_lib_path[:-4]), # /nvc is appended in wrapper
            "--entity={}".format(ctx.attr.name),
            "--library-dir-in-path={}".format(work_library_file.path),
            "--library-dir-out-path={}".format(out_dir.path),
            "--",
        ] + vpi_flags,
        tools = [analyzer_x, ctx.executable._script] + artifacts,
        # Only seems to work from bazel 6.0.0 on.
        #toolchain = _NVC_TOOLCHAIN_TYPE,
        progress_message = "Elaborating VHDL: {}.{} at {}".format(
            vhdl_provider.library_name, ctx.attr.name, work_library_file.path),
    )
    return  [
        VHDLLibraryProvider(
            libraries = vhdl_provider.libraries + all_libraries,
            entities = vhdl_provider.entities,
            library_name = library_name,
            library_dir = out_dir,
            vpi_plugins = vhdl_provider.vpi_plugins,
        ),
        ElaborateProvider(entity=ctx.attr.name),
        DefaultInfo(
            files=depset([out_dir]),
            runfiles=runfiles,
        ),
    ]

vhdl_elaborate = rule(
    doc = "Elaborates a VHDL design using NVC.",
    implementation = _vhdl_elaborate,
    attrs = {
        "library": attr.label(
            doc = "The `vhdl_library` target to elaborate.",
        ),
        "_script": attr.label(
            default = _NVC_WRAPPER,
            executable = True,
            cfg = "host",
            doc = "Wrapper script to run NVC.",
        ),
        "standard": attr.string(
            default = _VHDL_STANDARD_DEFAULT,
            doc = "The VHDL standard to use for elaboration (e.g., '2019').",
        ),
    },
    toolchains = [
        _NVC_TOOLCHAIN_TYPE
    ],
)
