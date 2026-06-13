load("//internal:utils.bzl", "get_single_file_from", "get_nvc_ld_library_path")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


load("//internal:vhdl_library.bzl", "vhdl_library")
load("//internal:vhdl_elaborate.bzl", "vhdl_elaborate")


def _vhdl_test(ctx):
    """
    Mostly copied from _vhdl_run
    """
    nvc_info = ctx.toolchains[_NVC_TOOLCHAIN_TYPE].nvc_info
    analyzer_x = nvc_info.analyzer.files.to_list()[0]
    analyzer = analyzer_x.short_path
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
    for name, path in all_libraries:
        if name != vhdl_provider.library_name and name not in seen:
            flag_libraries += [
                "--map={}:{}/{}".format(name, path.short_path, name)]
            deps_paths += [path]
            seen += [name]


    work_library_file = get_single_file_from(ctx.attr.entity)

    elaborate_provider = ctx.attr.entity[ElaborateProvider]

    runfiles = ctx.runfiles(
        files=[vhdl_provider.library_dir] + ctx.attr._script.files.to_list(),
        transitive_files=depset(artifacts + deps_paths))

    # Collect VPI plugins
    vpi_plugins = vhdl_provider.vpi_plugins.to_list()
    runfiles = runfiles.merge(ctx.runfiles(files = vpi_plugins))

    # Construct --load= flags for VPI plugins
    vpi_flags = ""
    if vpi_plugins:
        vpi_flags = "--load={}".format(",".join([p.short_path for p in vpi_plugins]))
    
    # Add user arguments
    extra_args = " ".join(ctx.attr.args)

    runfiles = runfiles.merge_all([ctx.attr._script[DefaultInfo].default_runfiles])
    inputs = deps_paths + [vhdl_provider.library_dir] + ([std_lib_dir] if hasattr(std_lib_dir, "path") else []) + artifacts + vpi_plugins
    i_runfiles = ctx.runfiles(files = inputs)

    tools = [analyzer_x, ctx.executable._script, ] + artifacts
    t_runfiles = ctx.runfiles(files = tools)

    runfiles.merge_all([t_runfiles, i_runfiles, ctx.attr._script[DefaultInfo].default_runfiles])
    
    # Calculate the runfiles prefix path
    # If the workspace is not _main (e.g. we are an external repo), we need to prefix the paths
    nvc_lib_path_for_wrapper = nvc_lib_path
    if nvc_lib_path.startswith("external/"):
        nvc_lib_path_for_wrapper = "../" + nvc_lib_path[9:]

    analyzer_for_wrapper = analyzer
    if analyzer.startswith("external/"):
        analyzer_for_wrapper = "../" + analyzer[9:]

    base_dir_for_wrapper = base_dir
    if base_dir.startswith("external/"):
        base_dir_for_wrapper = "../" + base_dir[9:]

    nvc_ld_library_path = get_nvc_ld_library_path(nvc_info, base_dir, ctx.configuration.default_shell_env)
    nvc_ld_library_path_for_wrapper = ":".join([
        ("../" + p[9:]) if p.startswith("external/") else p
        for p in nvc_ld_library_path.split(":")
    ])

    ctx.actions.expand_template(
        template = ctx.file._template,
        output = ctx.outputs.executable,
        substitutions = {
            "{{EXECUTABLE}}": "NVC_LD_LIBRARY_PATH=\"" + nvc_ld_library_path_for_wrapper + "\" LD_LIBRARY_PATH=\"" + base_dir_for_wrapper + "/lib/x86_64-linux-gnu\" " + ctx.executable._script.short_path,
            "{{VHDL_STANDARD}}": ctx.attr.standard,
            "{{ANALYZER}}": analyzer_for_wrapper,
            "{{LIBRARY_NAME}}": vhdl_provider.library_name,
            "{{LIBRARY_PATHS}}": " ".join(flag_libraries + ["-L", nvc_lib_path_for_wrapper]),
            "{{STDLIB_DIR}}": nvc_lib_path_for_wrapper[:-4],
            "{{ENTITY}}": elaborate_provider.entity,
            "{{LIB_DIR_IN_PATH}}": vhdl_provider.library_dir.short_path,
            "{{LIB_DIR_OUT_PATH}}": work_library_file.short_path,
            "{{WAVE_FILE}}": "{}.fst".format(ctx.attr.name),
            "{{VPI_FLAGS}}": vpi_flags,
            "{{EXTRA_ARGS}}": extra_args,
        },
    )
    return [DefaultInfo(runfiles=runfiles)]

_vhdl_internal_test = rule(
    doc = "Internal rule to execute a VHDL test using NVC.",
    test = True,
    implementation = _vhdl_test,
    attrs = {
        "entity": attr.label(
            doc = "The elaborated entity to test.",
        ),
        "deps": attr.label_list(
            default = [],
            doc = "Dependencies required for the test.",
        ),
        "_script": attr.label(
            default = _NVC_WRAPPER,
            executable = True,
            cfg = "host",
            doc = "Wrapper script to run NVC.",
        ),
        "standard": attr.string(
            default = _VHDL_STANDARD_DEFAULT,
            doc = "The VHDL standard to use (e.g., '2008', '2019').",
        ),
        "_template": attr.label(
            default = Label("//build/nvc:unittest.tpl.sh"),
            allow_single_file = True,
            doc = "Template for the test execution script.",
        ),
    },
    toolchains = [
      _NVC_TOOLCHAIN_TYPE,
    ],
)

def vhdl_test(name, srcs, deps,
    standard=_VHDL_STANDARD_DEFAULT, args=[], entity=None, entities=[]):
    """
    Defines a VHDL test.

    This macro combines `vhdl_library`, `vhdl_elaborate`, and internal test
    execution steps into a single logical target.

    Args:
        name: The name of the base test target.
        srcs: A list of VHDL source files (`.vhdl` or `.vhd`).
        deps: A list of `vhdl_library` targets that this test depends on.
        standard: The VHDL standard to use (e.g., "2008", "2019"). Defaults to "2019".
        args: A list of additional command-line arguments to pass to the NVC simulator.
        entity: A single entity to test.
        entities: A list of entities to test. If both `entity` and `entities` are provided, all are tested.
    """
    entity_list = []
    if entity:
        entity_list += [entity]
    entity_list += entities

    for entity in entity_list:
        # Strictly speaking this is not correct, since we're compiling the same
        # library twice.
        vhdl_library_name = "{name}_{entity}_lib".format(name=name,entity=entity)
        vhdl_library(
            name = vhdl_library_name,
            srcs = srcs,
            deps = deps,
            standard = standard,
        )
        e = "{entity}".format(entity=entity)
        vhdl_elaborate(
            name = e,
            library = ":{}".format(vhdl_library_name),
            standard = standard,
        )
        _vhdl_internal_test(
            name = "{name}_{entity}_test".format(name=name,entity=entity),
            entity = ":{}".format(e),
            args = args,
            standard = standard,
        )

