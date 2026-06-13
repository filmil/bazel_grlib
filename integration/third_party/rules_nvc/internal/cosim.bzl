# LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")
load("@rules_verilator//verilator:defs.bzl", "verilator_cc_library")
load("@rules_verilog//verilog:defs.bzl", vlog_library = "verilog_library")
load("//internal:vhdl_library.bzl", "vhdl_library")
load("//internal:vhdl_elaborate.bzl", "vhdl_elaborate")
load("//internal:vhdl_test.bzl", "vhdl_test")

def _verilator_json_impl(ctx):
    # Use the verilator toolchain
    verilator_toolchain = ctx.toolchains["@rules_verilator//verilator:toolchain_type"]
    verilator = verilator_toolchain.verilator

    json_out = ctx.actions.declare_file("{}_verilator.json".format(ctx.attr.name))

    # Verilator needs to find its includes. We pass VERILATOR_ROOT via environment,
    # computing it from the path to the executable.
    # The executable is typically at <VERILATOR_ROOT>/bin/verilator
    
    args = ctx.actions.args()
    args.add("--json-only")

    args.add("--top-module", ctx.attr.top_module)
        
    for k, v in ctx.attr.parameters.items():
        args.add("-G{}={}".format(k, v))

    args.add_all(ctx.files.srcs)

    # Need VERILATOR_ROOT to be correctly set so it finds includes.
    # It must point to the root directory where the "include/" folder is.
    
    ctx.actions.run_shell(
        outputs = [json_out],
        inputs = ctx.files.srcs + verilator_toolchain.all_files.to_list(),
        tools = [verilator],
        command = """
            # Find the waiver file to determine VERILATOR_ROOT
            WAIVER=$(find . -name verilated_std_waiver.vlt | head -n 1)
            export VERILATOR_ROOT=$(dirname $(dirname $WAIVER))
            
            {verilator} $@
            
            # Verilator --json-only outputs to an unspecified file. 
            # In latest versions it usually generates obj_dir/<top>.tree.json.
            find . -name "*.json"
            cp $(find . -name "*.tree.json" | head -n 1) {out}
        """.format(
            verilator = verilator.path,
            top = ctx.attr.top_module,
            out = json_out.path,
        ),
        arguments = [args],
        progress_message = "Generating Verilator JSON for {}".format(ctx.attr.name),
    )
    return [DefaultInfo(files = depset([json_out]))]

verilator_json = rule(
    implementation = _verilator_json_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True, mandatory = True),
        "top_module": attr.string(mandatory = True),
        "parameters": attr.string_dict(doc = "Verilog parameters (VHDL generics) to pass to the Verilator build"),
    },
    toolchains = ["@rules_verilator//verilator:toolchain_type"],
)
def _bridge_gen_impl(ctx):
    vhdl_out = ctx.actions.declare_file("{}_dut_proxy.vhdl".format(ctx.attr.name))
    hpp_out = ctx.actions.declare_file("{}_vpi_bindings.hpp".format(ctx.attr.name))
    
    outputs = [vhdl_out, hpp_out]
    arguments = [
        ctx.files.json_src[0].path,
        vhdl_out.path,
        hpp_out.path,
        ctx.attr.top_module,
    ]

    vhdl_arch_out = None
    if ctx.attr.separate_entity_arch:
        vhdl_arch_out = ctx.actions.declare_file("{}_dut_proxy_arch.vhdl".format(ctx.attr.name))
        outputs.append(vhdl_arch_out)
        arguments.append(vhdl_arch_out.path)

    ctx.actions.run(
        outputs = outputs,
        inputs = ctx.files.json_src,
        executable = ctx.executable._generator,
        arguments = arguments,
        progress_message = "Generating Co-Simulation Bridge for {}".format(ctx.attr.name),
    )
    
    vhdl_depset = depset([vhdl_out])
    vhdl_arch_depset = depset([vhdl_arch_out]) if vhdl_arch_out else depset()

    return [
        DefaultInfo(files = depset(outputs)),
        OutputGroupInfo(
            vhdl = vhdl_depset,
            vhdl_arch = vhdl_arch_depset,
            hpp = depset([hpp_out]),
        )
    ]

bridge_gen = rule(
    implementation = _bridge_gen_impl,
    attrs = {
        "json_src": attr.label(allow_files = True, mandatory = True),
        "top_module": attr.string(mandatory = True),
        "separate_entity_arch": attr.bool(default = False),
        "_generator": attr.label(
            default = Label("//internal/cosim:generate_bridge"),
            executable = True,
            cfg = "exec",
        ),
    },
)

def _bridge_vhdl_impl(ctx):
    if ctx.attr.arch_only:
        return [DefaultInfo(files = ctx.attr.bridge[OutputGroupInfo].vhdl_arch)]
    return [DefaultInfo(files = ctx.attr.bridge[OutputGroupInfo].vhdl)]

bridge_vhdl = rule(
    implementation = _bridge_vhdl_impl,
    attrs = {
        "bridge": attr.label(mandatory = True),
        "arch_only": attr.bool(default = False),
    },
)

def _bridge_hpp_impl(ctx):
    return [DefaultInfo(files = ctx.attr.bridge[OutputGroupInfo].hpp)]

bridge_hpp = rule(
    implementation = _bridge_hpp_impl,
    attrs = {
        "bridge": attr.label(mandatory = True),
    },
)

def nvc_verilator_cosim(name, srcs, top_module, path_prefix, parameters = {}, separate_entity_arch = False, standard = "2019"):
    """
    Macro that encapsulates NVC and Verilator co-simulation bindings generation.
    Generates the VHDL proxy and C++ bindings for the given Verilog top module.
    The resulting VHDL file can be used as a dependency in a `vhdl_library`.
    
    Args:
        name: Name of the generated rules.
        srcs: Verilog source files.
        top_module: The top module to compile.
        path_prefix: The exact VHDL hierarchical path where this component will
                     be instantiated in the testbench. NVC requires this absolute
                     string (e.g., ":top_tb:dut_inst") to natively resolve the 
                     VHPI handles at runtime since dynamic hierarchy discovery
                     from within VHPIDIRECT is heavily restricted. The user must 
                     manually align this string with their VHDL architecture names 
                     and component instantiation labels.
        parameters: Verilog parameters (VHDL generics) to pass to Verilator.
        separate_entity_arch: If True, generate separate entity and architecture files.
                              Creates an additional target <name>.archonly.
        standard: The VHDL standard to use for co-simulation bridge generation.
    """
    json_name = "{}_json".format(name)
    verilator_json(
        name = json_name,
        srcs = srcs,
        top_module = top_module,
        parameters = parameters,
    )

    bridge_name = "{}_bridge".format(name)
    bridge_gen(
        name = bridge_name,
        json_src = ":{}".format(json_name),
        top_module = top_module,
        separate_entity_arch = separate_entity_arch,
    )

    vhdl_entity_target = "{}.vhdl".format(name)
    bridge_vhdl(
        name = vhdl_entity_target,
        bridge = ":{}".format(bridge_name),
    )

    hpp_name = "{}_hpp".format(name)
    bridge_hpp(
        name = hpp_name,
        bridge = ":{}".format(bridge_name),
    )

    vlog_lib_name = "{}_vlog".format(name)
    vlog_library(
        name = vlog_lib_name,
        srcs = srcs,
    )

    verilated_lib_name = "{}_verilated".format(name)
    verilator_cc_library(
        name = verilated_lib_name,
        module = ":{}".format(vlog_lib_name),
        module_top = top_module,
        vopts = ["-G{}={}".format(k, v) for k, v in parameters.items()],
    )

    cc_name = "{}.vpi".format(name)
    cc_binary(
        name = cc_name,
        srcs = [
            Label("//internal/cosim:vpi_wrapper.cpp"), 
            Label("//third_party/ieee:vhpi_user.h"), 
            ":{}".format(hpp_name)
        ],
        copts = [
            "-I$(GENDIR)/" + native.package_name(), 
            "-I$(BINDIR)/" + native.package_name(),
            "-Iexternal/rules_nvc+/third_party/ieee",
            "-DVPI_BINDINGS_HEADER=\\\"{}_bridge_vpi_bindings.hpp\\\"".format(name),
            "-DVERILATOR_STEP_CALL=verilator_step_call_{}".format(top_module),
            "-DPATH_PREFIX=\\\"{}\\\"".format(path_prefix)
        ],
        linkshared = True,
        deps = [
            ":{}".format(verilated_lib_name),
        ], 
    )
    
    vhdl_srcs = [":{}".format(vhdl_entity_target)]
    
    if separate_entity_arch:
        vhdl_arch_target = "{}.arch_vhdl".format(name)
        bridge_vhdl(
            name = vhdl_arch_target,
            bridge = ":{}".format(bridge_name),
            arch_only = True,
        )
        vhdl_srcs.append(":{}".format(vhdl_arch_target))
        
        native.alias(
            name = name + ".archonly",
            actual = vhdl_arch_target,
        )

    vhdl_library(
        name = name,
        srcs = vhdl_srcs,
        vpi_plugins = [":{}".format(cc_name)],
        standard = standard,
    )
