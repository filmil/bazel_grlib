load("//internal:utils.bzl", "get_single_file_from")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


def _produce_waveform(ctx):
    data_files = []
    for target in ctx.attr.data:
        for file in target.files.to_list():
            data_files += [file]

    ext = "vcd"
    if ctx.attr.use_fst:
        ext = "fst"

    output_file = ctx.actions.declare_file("{}.{}".format(ctx.attr.name, ext))
    runfiles = ctx.runfiles(files=[output_file] + data_files)
    sim = ctx.executable.simulation
    ctx.actions.run(
        outputs = [output_file],
        inputs = data_files,
        executable = sim.path,
        arguments = [
            output_file.path,
        ] + ctx.attr.args,
        tools = [
            sim,
        ],
    )
    return [
        DefaultInfo(
            files=depset([output_file]),
            runfiles=runfiles,
        ),
    ]


produce_waveform = rule(
    doc = "Produces a waveform file (VCD) from a VHDL simulation run.",
    implementation = _produce_waveform,
    attrs = {
        "simulation": attr.label(
            executable = True,
            cfg = "host",
            doc = "The simulation target (`vhdl_run`) to execute.",
        ),
        "data": attr.label_list(
            doc = "Data files required for the simulation.",
        ),
        "args": attr.string_list(
            doc = "Additional command-line arguments to pass to the simulation.",
        ),
        "use_fst": attr.bool(
            default = False,
            doc = "A boolean indicating whether to expect an FST file instead of VCD file. Defaults to `False`.",
        ),
    },
)

