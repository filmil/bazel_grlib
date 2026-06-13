# NVC and Verilator Co-Simulation Bridge Report

This document outlines the technical challenges and solutions discovered while
building the co-simulation bridge between NVC (VHDL) and Verilator (C++).

## 1. VPI to VHPI Migration
**Problem:** Initially, the bridge relied on standard `VPI` functions
(`vpi_handle_by_name`, `vpi_put_value`, `vpi_get_value`). However, NVC's
implementation of VPI is extremely limited and primarily designed to support
static system tasks (`$display`). Specifically, `vpi_handle_by_name` is
entirely unimplemented in NVC, and iterating the hierarchy dynamically during a
`VHPIDIRECT` execution returns `NULL`.
**Solution:** We migrated the bridge to use `VHPI` (`vhpi_user.h`), which is
fully native to NVC. Functions like `vhpi_handle_by_name` and `vhpi_get_value`
are dynamically exported and function correctly in the execution context of a
VHDL testbench.

## 2. Dynamic Symbol Resolution vs Verilator Stubs
**Problem:** Verilator statically links its own dummy VPI/VHPI implementations
(`verilated_vpi.cpp`) into the compiled model. When our C++ bridge called
`vhpi_handle_by_name`, the linker bound it to Verilator's stub (which simply
returns `nullptr`), causing NVC to crash or fail to bind signals.
**Solution:** We bypassed static linkage by leveraging `dlsym(RTLD_DEFAULT,
"vhpi_get_value")` inside `vpi_wrapper.cpp` to dynamically load the *actual*
VHPI pointers from the running NVC simulator binary at runtime. However, we
later discovered that NVC natively exports `vhpi_handle_by_name` and
`vhpi_put_value` globally, so compiling it purely against the standard
`<vhpi_user.h>` (without `dlsym`) works reliably.

## 3. Disambiguating Multiple C++ Bridges
**Problem:** In `demo_test`, we needed to instantiate both `adder` and `counter`
components. When attempting to link two separate C++ bridges into the same
testbench, the static `verilator_step_call` C function symbol collided.
**Solution:** We introduced the `VERILATOR_STEP_CALL` C-preprocessor macro to
append the module's name to the VHPIDIRECT function (e.g.,
`verilator_step_call_adder` and `verilator_step_call_counter`). The generated
VHDL proxy dynamically binds to this uniquely named symbol via `attribute
foreign of step_verilator : procedure is "VHPIDIRECT
verilator_step_call_{{.Name}}";`.

## 4. Scalar vs Vector Logic Value Formatting
**Problem:** When transferring `std_logic_vector` sizes > 1, mapping to the
`vhpiIntVal` format silently truncated values in NVC.
**Solution:** We updated `generate_bridge.go` to explicitly query the bit-width
of each handle using `vhpi_get(vhpiSizeP, net_handle)`. 
- If `size == 1`, we use `vhpiLogicVal` and map to `val.value.enumv`.
- If `size > 1`, we use `vhpiLogicVecVal` and construct a `vhpiEnumT` array to
  pass the signal value bit-by-bit (treating `3` as `1` and `2` as `0`).

## 5. Endianness / Bit Ordering
**Problem:** VHDL models expected data where index `0` was the Most Significant
Bit (MSB), whereas naive integer casting mapped it to the Least Significant
Bit.
**Solution:** We updated the `for` loops within `sync_inputs` and `sync_outputs`
to use `(1ULL << (size - 1 - i))` to explicitly align Verilator's scalar C++
types into VHDL's `vhpiEnumT` arrays properly.

## 6. Verilator JSON Parameter Passing
**Problem:** Parameterized Verilog modules (`adder #(WIDTH=16)`) were not being
passed properly to `rules_verilator`, causing the C++ model to compile as 8-bit
while the JSON/VHDL interface exposed 16-bit ports.
**Solution:** We updated the `cosim.bzl` macro to explicitly forward the
`parameters = {}` dictionary downstream as `vopts = ["-G{}={}".format(k, v)]`
to `verilator_cc_library`.

## 8. Separate VHDL Entity and Architecture Generation
**Problem:** In some use cases, developers may want to manually write their VHDL
`entity` declaration for the co-simulated module. This is useful when the
auto-generated entity has generic names, missing custom attributes, or non-ideal
port mapping. However, the `generate_bridge` code tightly coupled both the
`entity` and the `architecture` in a single proxy file.
**Solution:** We added a `separate_entity_arch` parameter to the macro and a
5th argument to the `generate_bridge.go` tool. When enabled, it outputs two
files (one for the entity, one for the architecture) and provides an `.archonly`
alias target for Bazel. This enables a user to write their own `entity` and
merely link in the Verilator C++ bindings and VPI plugin using the
`vpi_plugins` and `deps` functionality.
