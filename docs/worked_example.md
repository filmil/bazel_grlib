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
