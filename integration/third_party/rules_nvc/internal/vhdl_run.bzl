load("//internal:utils.bzl", "get_single_file_from", "get_nvc_deps", "get_nvc_ld_library_path")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _vhdl_run(ctx):
    nvc_info = ctx.toolchains[_NVC_TOOLCHAIN_TYPE].nvc_info
    nvc_deps = get_nvc_deps(nvc_info)
    analyzer_x = nvc_info.analyzer.files.to_list()[0]
    analyzer = analyzer_x.path
    library_name = ctx.attr.name

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


    vhdl_provider = ctx.attr.entity[VHDLLibraryProvider]

    all_libraries = []
    flag_libraries = []
    deps_paths = []
    seen = []
    all_libraries += vhdl_provider.libraries
    for name, path in vhdl_provider.libraries:
        if name != vhdl_provider.library_name and name not in seen:
            flag_libraries += [
                "--map={}:{}/{}".format(name, path.path, name)]
            deps_paths += [path]
            seen += [name]


    work_library_file = get_single_file_from(ctx.attr.entity)

    ext = "vcd"
    if ctx.attr.use_fst:
        ext = "fst"

    wave_file = ctx.actions.declare_file(
        "{}.{}".format(ctx.attr.name, ext))

    format = []
    if ctx.attr.use_fst:
        format = [ "--format=fst" ]
    elif ctx.attr.use_vcd:
        format = [ "--format=vcd" ]

    elaborate_provider = ctx.attr.entity[ElaborateProvider]

    vpi_plugins = vhdl_provider.vpi_plugins.to_list()
    vpi_flags = []
    for p in vpi_plugins:
        vpi_flags += ["-m", p.path]

    runfiles = ctx.runfiles(files = [wave_file] + vpi_plugins)
    ctx.actions.run(
        outputs = [wave_file],
        inputs = depset(direct = deps_paths + [vhdl_provider.library_dir] + ([std_lib_dir] if hasattr(std_lib_dir, "path") else []) + artifacts + nvc_deps + vpi_plugins).to_list(),
        executable = ctx.executable._script.path,
        env = {
            "NVC_LD_LIBRARY_PATH": get_nvc_ld_library_path(nvc_info, base_dir, ctx.configuration.default_shell_env),
        },
        arguments = [
            "-cmd=-r",
            "--vhdl-standard={}".format(ctx.attr.standard),
            "--nvc-binary-path={}".format(analyzer),
            "--library-name={}".format(vhdl_provider.library_name),
            "--library-paths={}".format(" ".join(flag_libraries + ["-L", nvc_lib_path])),
            "--stdlib-dir={}".format(nvc_lib_path[:-4]),
            "--entity={}".format(elaborate_provider.entity),
            "--library-dir-in-path={}".format(work_library_file.path),
            "--library-dir-out-path={}".format(work_library_file.path),
            "--",
        ] + vpi_flags + ctx.attr.args + [
            "--wave={}".format(wave_file.path),
        ] + format,
        tools = [analyzer_x, ctx.executable._script] + artifacts,
        # Only seems to work from bazel 6.0.0 on.
        #toolchain = _NVC_TOOLCHAIN_TYPE,
        progress_message = "Simulating VHDL: {}.{} at {}".format(
            vhdl_provider.library_name,
            elaborate_provider.entity,
            work_library_file.path),
    )
    return [
        DefaultInfo(
            files=depset([wave_file]),
            runfiles=runfiles,
        ),
    ]

vhdl_run = rule(
    doc = "Simulates an elaborated VHDL design using NVC.",
    implementation = _vhdl_run,
    attrs = {
        "entity": attr.label(
            doc = "The elaborated VHDL entity to simulate. This should be a `vhdl_elaborate` target.",
        ),
        "deps": attr.label_list(
            default = [],
            doc = "A list of other `vhdl_library` targets that this simulation depends on.",
        ),
        "_script": attr.label(
            default = _NVC_WRAPPER,
            executable = True,
            cfg = "host",
            doc = "Wrapper script to run NVC.",
        ),
        "standard": attr.string(
            default = _VHDL_STANDARD_DEFAULT,
            doc = "The VHDL standard to use for simulation. Defaults to '2019'.",
        ),
        "use_vcd": attr.bool(
            default = True,
            doc = "A boolean indicating whether to generate a VCD (Value Change Dump) file for waveform viewing. Defaults to `True`.",
        ),
        "use_fst": attr.bool(
            default = False,
            doc = "A boolean indicating whether to generate a FST file for waveform viewing. Defaults to `False`. Takes precedence over `use_vcd`.",
        ),
        "args": attr.string_list(
            doc = "A list of added command line args to use",
        ),
    },
    toolchains = [
      _NVC_TOOLCHAIN_TYPE,
    ],
)
