# Co-Simulation Usage Example

This guide demonstrates how to use the NVC and Verilator co-simulation bridge to
instantiate a parameterized Verilog module within a VHDL testbench.

## 1. Verilog Module (`adder.v`)

First, define your Verilog module. Here we have a simple synchronous adder with
a configurable `WIDTH` parameter.

```verilog
module adder #(
    parameter WIDTH = 8
) (
    input wire clk,
    input wire [WIDTH-1:0] a,
    input wire [WIDTH-1:0] b,
    output reg [WIDTH:0] sum
);
    always @(posedge clk) begin
        sum <= a + b;
    end
endmodule
```

## 2. Bazel Build Configuration (`BUILD.bazel`)

Use the `nvc_verilator_cosim` macro to generate the co-simulation bridge. You
can override the default Verilog parameters using the `parameters` attribute.
The `path_prefix` attribute must match the hierarchical path of the component
instantiation in your VHDL testbench.

```starlark
load("@rules_nvc//internal:cosim.bzl", "nvc_verilator_cosim")
load("@rules_nvc//nvc:rules.bzl", "vhdl_test")

nvc_verilator_cosim(
    name = "cosim_bridge",
    srcs = ["adder.v"],
    top_modules = ["adder"],
    parameters = {
        "WIDTH": "16",
    },
    path_prefix = ":top_tb:dut_inst",
)

vhdl_test(
    name = "cosim_test",
    srcs = ["top_tb.vhdl"],
    deps = [
        ":cosim_bridge",
    ],
    entity = "top_tb",
)
```

### Separate Entity and Architecture Generation

If you prefer to define the VHDL `entity` manually (e.g., to add custom generics
or modify port types), you can set `separate_entity_arch = True`. This exposes
an `.archonly` alias target which provides only the generated architecture and
C++ bindings.

```starlark
# 1. Generate only the architecture
nvc_verilator_cosim(
    name = "cosim_bridge_gen",
    srcs = ["adder.v"],
    top_module = "adder",
    path_prefix = ":top_tb:dut_inst",
    separate_entity_arch = True,
)

# 2. Combine your manual entity with the generated architecture
vhdl_library(
    name = "cosim_bridge",
    srcs = [
        "adder_entity.vhdl", # Your manually written entity
    ],
    deps = [
        ":cosim_bridge_gen.archonly",
    ],
)
```

## 3. VHDL Testbench (`top_tb.vhdl`)

In your VHDL testbench, instantiate the generated proxy entity directly from the
compiled bridge library. The `INSTANCE_ID` generic is required by the bridge
to track multiple instances of the same module in the C++ state.

```vhdl
library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

-- Make the generated bridge library visible. The library name matches the
-- name of the nvc_verilator_cosim target.
library cosim_bridge;

entity top_tb is
end entity;

architecture sim of top_tb is
  signal clk : std_logic := '0';
  signal a   : std_logic_vector(15 downto 0) := x"1234";
  signal b   : std_logic_vector(15 downto 0) := x"4321";
  signal sum : std_logic_vector(16 downto 0);
begin
  clk <= not clk after 5 ns;

  -- Instantiate the generated proxy entity directly from the library
  -- Note: The path ":top_tb:dut_inst" matches the path_prefix in BUILD.bazel
  dut_inst : entity cosim_bridge.adder
    generic map (
      INSTANCE_ID => 1
    )
    port map (
      clk => clk,
      a   => a,
      b   => b,
      sum => sum
    );

  process
  begin
    wait for 20 ns;
    
    -- Verify that Verilator successfully computed 0x1234 + 0x4321 = 0x5555
    -- using the overridden WIDTH=16 parameter.
    assert sum = '0' & x"5555" report "Sum incorrect!" severity failure;
    
    report "Simulation finished successfully";
    std.env.finish;
  end process;
end architecture;
```
