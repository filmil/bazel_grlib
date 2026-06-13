# Worked Example: Configuring a LEON3 System

This guide walks through the process of customizing a GRLIB hardware design using the Bazel/Kconfig system.

## Scenario

We want to configure a LEON3-based system with the following specific hardware parameters:

1.  **Register Windows:** Increase to **16** (from the default).
2.  **Floating Point Unit:** Ensure the FPU is **enabled**.
3.  **AHB Data Width:** Set the internal bus width to **64 bits**.

## Step 1: Identify the Kconfig Symbols

First, we locate the namespaced symbols for our desired parameters in `docs/example.bazelrc`:

*   **Register Windows:** `@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS`
*   **FPU Enable:** `@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_FPU_ENABLE`
*   **AHB Data Width:** `@grlib_config//:CONFIG_AHBDW` (this is a base item, not namespaced)

## Step 2: Set the Active Design Prefix

To make our `BIN_TKCONFIG` settings the "global" ones used by the IP cores, we must set the `ACTIVE_DESIGN_PREFIX` to match the namespace we are using.

```bash
--@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=BIN_TKCONFIG_CONFIG
```

### What does "Promotion" mean?

Because GRLIB is a large library with many components, different parts of the system might define the same configuration parameter (e.g., `IU_NWINDOWS`). To allow building complex SoCs where different cores might have different settings, this build system **namespaces** every parameter based on its file path:

*   `IU_NWINDOWS` in `lib/gaisler/leon3/leon3.in` becomes `LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS`.

However, the VHDL IP core code itself is generic and expects constants with standard names like `CFG_IU_NWINDOWS` inside the `grlib.config` package.

**Promotion** is the process where the build system takes the namespaced symbols belonging to your **Active Design** and maps them to these generic names. When you set `ACTIVE_DESIGN_PREFIX=BIN_TKCONFIG_CONFIG`:

1.  The generator finds `CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS`.
2.  It "promotes" it to the generic `CFG_IU_NWINDOWS` in the final VHDL package.
3.  The LEON3 IP core, which uses `CFG_IU_NWINDOWS`, now automatically uses your specific value.

This allows you to maintain a library of 1,100+ unique parameters while still providing the standard interface the hardware logic expects.

### Handling Multiple Components with Same Setting Names

What if you have two different processor cores in your SoC (e.g., a LEON3 and a LEON4), both of which have a parameter named `IU_NWINDOWS`?

GRLIB's legacy system would force you to use the same value for both. However, this Bazel build system's **namespacing** allows you to configure them independently:

1.  **LEON3 Setting:** `@grlib_config//:CONFIG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS`
2.  **LEON4 Setting:** `@grlib_config//:CONFIG_LIB_GAISLER_LEON4_LEON4_IU_NWINDOWS`

You can set these to different values in your `.bazelrc`:

```text
build --@grlib_config//:CONFIG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS=8
build --@grlib_config//:CONFIG_LIB_GAISLER_LEON4_LEON4_IU_NWINDOWS=16
```

### Instantiating Multiple Different Cores in VHDL

Since all namespaced symbols are exported as constants in the `grlib.config` package, you can pass them explicitly to the IP core generics during instantiation.

This bypasses the global "promoted" defaults and allows multiple core variants to coexist in the same SoC:

```vhdl
library grlib;
use grlib.config.all;

-- ... inside your architecture ...

-- 1. A LEON3 with 8 windows
cpu0: entity gaisler.leon3
  generic map (
    nwindows => CFG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS,
    -- ... other generics ...
  )
  port map ( ... );

-- 2. A LEON4 with 16 windows
cpu1: entity gaisler.leon4
  generic map (
    nwindows => CFG_LIB_GAISLER_LEON4_LEON4_IU_NWINDOWS,
    -- ... other generics ...
  )
  port map ( ... );
```

Note that the "Promotion" logic (via `ACTIVE_DESIGN_PREFIX`) is intended for top-level SoC designs where you want to apply a specific set of parameters to all generic components. For fine-grained control of individual IP cores, use these namespaced symbols directly.

## Step 3: Run the Build

You can apply these settings directly on the command line:

```bash
bazel build \
  --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=BIN_TKCONFIG_CONFIG \
  --@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS=16 \
  --@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_FPU_ENABLE=True \
  --@grlib_config//:CONFIG_AHBDW=64 \
  @grlib_srcs//:grlib
```

## Step 4: Verify the Configuration

When you run the build above, Bazel programmatically generates a VHDL package. You can inspect the results in the Bazel output directory:

```bash
# Locate and read the generated VHDL config
cat bazel-bin/third_party/grlib/config.vhd
```

In the output, you will see your parameters promoted to the generic names expected by GRLIB:

```vhdl
-- ...
package config is
  -- Namespaced symbols
  constant CFG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS : integer := 16;
  -- ...

  -- Promoted symbols for prefix: BIN_TKCONFIG_CONFIG
  constant CFG_IU_NWINDOWS : integer := 16;
  constant CFG_FPU_ENABLE : integer := 1;
  constant CFG_AHBDW : integer := 64;
  -- ...
end;
```

## Step 5: Persisting Settings

To avoid typing long flags every time, add them to your `user.bazelrc` at the root of the project:

```text
# user.bazelrc
build --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=BIN_TKCONFIG_CONFIG
build --@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS=16
build --@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_FPU_ENABLE=True
build --@grlib_config//:CONFIG_AHBDW=64
```

Now, a simple `bazel build @grlib_srcs//:grlib` will always use your custom hardware configuration.

## Advanced: Building Multiple Configurations via Transitions

If you need to build the same library with different settings in the same workspace (e.g., a "Lite" vs "Full" version of a core), you can use **Bazel Transitions**.

### 1. Define the Transition (`transitions.bzl`)

Create a Starlark file to define how the configuration should be modified for specific targets:

```python
# transitions.bzl

def _leon3_high_perf_transition_impl(settings, attr):
    # This transition overrides specific Kconfig settings
    return {
        "@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX": "BIN_TKCONFIG_CONFIG",
        "@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS": "16",
        "@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_FPU_ENABLE": "True",
    }

leon3_high_perf_transition = transition(
    implementation = _leon3_high_perf_transition_impl,
    inputs = [],
    outputs = [
        "@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX",
        "@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_IU_NWINDOWS",
        "@grlib_config//:CONFIG_BIN_TKCONFIG_CONFIG_FPU_ENABLE",
    ],
)

def _configured_library_impl(ctx):
    # Forward the providers from the underlying library built in the new configuration
    return [ctx.attr.library[0][DefaultInfo]]

configured_library = rule(
    implementation = _configured_library_impl,
    attrs = {
        "library": attr.label(cfg = leon3_high_perf_transition),
        "_allowlist_function_transition": attr.label(
            default = "@bazel_tools//tools/allowlists/function_transition_allowlist",
        ),
    },
)
```

### 2. Use the Rule in your BUILD file

Now you can define targets that represent specific hardware variants:

```python
# BUILD.bazel
load(":transitions.bzl", "configured_library")

configured_library(
    name = "leon3_high_perf",
    library = "@grlib_srcs//:grlib",
)
```

### 3. Build the Variant

Building this target will automatically apply the configuration changes defined in the transition, regardless of what is set on your command line or in `.bazelrc`:

```bash
bazel build //:leon3_high_perf
```

## Worked Example: Configuring and Simulating NOEL-V (64-bit)

This example shows how to configure the high-performance NOEL-V RISC-V core for a 64-bit simulation.

### 1. Identify the Configuration Flags

To enable NOEL-V in 64-bit mode, we need to promote its specific namespace and set the data width:

*   **Active Design Prefix:** `LIB_GAISLER_NOELV_NOELV`
*   **Processor Enable:** `CONFIG_LIB_GAISLER_NOELV_NOELV_NOELV=True`
*   **Data Width (XLEN):** `CONFIG_LIB_GAISLER_NOELV_NOELV_NOELV_XLEN=64`

### 2. Build for Simulation

Use the following command to build the core libraries with these 64-bit parameters:

```bash
bazel build \
  --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=LIB_GAISLER_NOELV_NOELV \
  --@grlib_config//:CONFIG_LIB_GAISLER_NOELV_NOELV_NOELV=True \
  --@grlib_config//:CONFIG_LIB_GAISLER_NOELV_NOELV_NOELV_XLEN=64 \
  @grlib_srcs//:grlib @grlib_srcs//:techmap
```

### 3. Verification

The build will generate a namespaced and promoted `config.vhd` package. You can verify that `CFG_NOELV_XLEN` is set to 64:

```bash
grep "constant CFG_NOELV_XLEN" bazel-bin/third_party/grlib/config.vhd
```
