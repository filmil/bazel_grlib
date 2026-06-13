<!-- LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4 -->

# Simulation vs. Synthesis Flow

This document explains how to integrate the NVC-Verilator co-simulation bridge
into a project that targets both high-performance simulation and hardware
synthesis.

## The Challenge

The co-simulation bridge generates a VHDL **proxy** that replaces your actual
Verilog/SystemVerilog RTL during simulation. This proxy forwards all port
values to a C++ model compiled by Verilator. 

While this is excellent for simulation speed and C++-based verification, the 
proxy and its associated VPI plugin are **not synthesizable**. Hardware 
synthesis tools (like Vivado, Quartus, or Yosys) require the original 
Verilog/SystemVerilog source files.

## Recommended Project Structure

To support both flows, it is recommended to use VHDL **configurations** or 
library-based switching to swap between the co-simulation proxy and the real 
RTL.

### 1. Dual-Library Approach

In this approach, you keep the generated co-simulation bridge in one library 
(e.g., `cosim_lib`) and your actual synthesizable RTL (or its VHDL wrapper) in 
another (e.g., `work` or `rtl_lib`).

#### Simulation (Bazel)

In your `BUILD.bazel`, the `vhdl_test` target depends on the 
`nvc_verilator_cosim` target.

```starlark
nvc_verilator_cosim(
    name = "dut_cosim",
    srcs = ["dut.v"],
    top_module = "dut",
    path_prefix = ":top_tb:dut_inst",
)

vhdl_test(
    name = "sim_test",
    srcs = ["top_tb.vhdl"],
    deps = [":dut_cosim"],
    entity = "top_tb",
)
```

In `top_tb.vhdl`, you instantiate the entity from the `dut_cosim` library:

```vhdl
library dut_cosim;
...
dut_inst : entity dut_cosim.dut
  generic map ( INSTANCE_ID => 1 )
  port map ( ... );
```

#### Synthesis (Hardware Flow)

For synthesis, you ignore the `nvc_verilator_cosim` target and instead provide 
the actual `dut.v` to your synthesis tool. Your synthesis tool will see the 
`dut` module directly.

If your top-level is VHDL, you would instantiate the `dut` component normally, 
and the synthesis tool will link it to the compiled Verilog module.

### 2. Using VHDL Configurations

You can define a single testbench that instantiates a component, and then use 
different VHDL configurations to bind that component to either the co-simulation 
proxy or the real RTL.

```vhdl
-- Testbench
architecture arch of top_tb is
  component dut is
    port ( ... );
  end component;
begin
  dut_inst : dut port map ( ... );
end architecture;

-- Configuration for co-simulation
configuration cfg_cosim of top_tb is
  for arch
    for dut_inst : dut
      use entity cosim_lib.dut(proxy)
        generic map ( INSTANCE_ID => 1 );
    end for;
  end for;
end configuration;

-- Configuration for synthesis/gate-level sim
configuration cfg_synth of top_tb is
  for arch
    for dut_inst : dut
      use entity work.dut(rtl);
    end for;
  end for;
end configuration;
```

## Automating the Switch in Bazel

You can use Bazel `select()` statements to automatically switch between 
co-simulation and real RTL based on a command-line flag or a build configuration.

```starlark
# Define a config_setting for synthesis
config_setting(
    name = "synthesis_mode",
    values = {"define": "mode=synthesis"},
)

vhdl_library(
    name = "dut_lib",
    srcs = select({
        ":synthesis_mode": ["dut_rtl_wrapper.vhdl"],
        "//conditions:default": [":dut_cosim"],
    }),
)
```

## Summary of Differences

| Feature | Co-Simulation (NVC + Verilator) | Synthesis / Gate-Level |
| :--- | :--- | :--- |
| **Tool** | NVC Simulator | Vivado, Quartus, Yosys, etc. |
| **Execution** | C++ Model (Verilated) | Hardware Gates / FPGA Bits |
| **Speed** | Extremely Fast | N/A (Hardware execution) |
| **Accuracy** | Cycle-accurate RTL behavior | Bit-accurate / Timing-accurate |
| **Synthesizable** | **No** | **Yes** |
| **Debug** | GDB, C++ Printfs, FST/VCD | ChipScope, Oscilloscope, Post-Synth Sim |
