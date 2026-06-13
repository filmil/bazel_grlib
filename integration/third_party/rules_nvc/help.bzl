
def _help_impl(ctx):
    out = ctx.actions.declare_file("help.txt")
    tc = ctx.toolchains["@rules_verilator//verilator:toolchain_type"]
    ctx.actions.run_shell(
        outputs = [out],
        inputs = [],
        tools = [tc.verilator],
        command = "{} --help > {} 2>&1 || true".format(tc.verilator.path.replace("verilator", "bin/verilator"), out.path),
    )
    return [DefaultInfo(files = depset([out]))]

verilator_help = rule(
    implementation = _help_impl,
    toolchains = ["@rules_verilator//verilator:toolchain_type"],
)
