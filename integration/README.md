# Reusing bazel_grlib in Your Project

This directory demonstrates how to use `bazel_grlib` as a dependency in an external Bazel workspace.

## 1. Bzlmod Setup

Add the following to your `MODULE.bazel` file:

```python
bazel_dep(name = "grlib", version = "1.0.0")
# If using a local clone for development:
local_path_override(
    module_name = "grlib",
    path = "path/to/bazel_grlib",
)

# You also need rules_nvc and rules_kconfig
bazel_dep(name = "rules_nvc", version = "4.2.4")
bazel_dep(name = "rules_kconfig", version = "0.4.0")
```

## 2. Configuration

`bazel_grlib` uses a namespaced Kconfig system. You can set hardware parameters for any GRLIB component using Bazel build settings.

### Using Global Settings

To use the global configuration package provided by `grlib`:

1.  In your `BUILD` file, depend on `@grlib//:grlib`.
2.  Set the `ACTIVE_DESIGN_PREFIX` to promote a specific component's settings.

```bash
bazel build --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=LIB_GAISLER_NOELV_NOELV //:my_target
```

### Reference Settings

Check `@grlib//docs:example.bazelrc` in the GRLIB repository for a full list of available namespaced parameters.

## 3. Building and Simulating

Instantiate GRLIB components in your VHDL code as usual. In your `BUILD.bazel`, use `vhdl_library` or `vhdl_test` from `rules_nvc`:

```python
load("@rules_nvc//nvc:rules.bzl", "vhdl_test")

vhdl_test(
    name = "my_design_sim",
    srcs = ["my_tb.vhd"],
    deps = [
        "@grlib//:grlib",
        "@grlib//:gaisler",
        "@grlib//:techmap",
    ],
)
```

## 4. Run the Simulation

```bash
bazel test //:my_design_sim
```
