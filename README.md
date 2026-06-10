# GRLIB Bazel Repository

![Build Status](https://github.com/filmil/bazel_grlib/actions/workflows/build.yml/badge.svg)

**Note:** The repository name is `bazel_grlib`, while the Bazel module name is `grlib`.

This repository contains a Bazel-based build system for GRLIB VHDL designs.

## Configuration

This repository uses `rules_kconfig` to manage hardware configuration. Every GRLIB component's configuration is namespaced to avoid collisions.

### Namespacing

Symbols are prefixed with their GRLIB path. For example, `IU_NWINDOWS` defined in `lib/gaisler/leon3/leon3.in` becomes:

`@grlib_config//:CONFIG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS`

### Promoting a Design Configuration

To make a specific design's configuration the "active" one (mapping namespaced symbols like `DESIGNS_LEON3MP_IU_NWINDOWS` back to the generic `CFG_IU_NWINDOWS` expected by VHDL), set the `ACTIVE_DESIGN_PREFIX`:

```bash
# Example: Activate the LEON3MP design configuration
bazel build --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=DESIGNS_LEON3MP //...
```

### Setting individual values

You can override any namespaced value via the command line:

```bash
bazel build --@grlib_config//:CONFIG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS=16 //...
```

### Reference

See `docs/worked_example.md` for a step-by-step guide on customizing your build, and `docs/example.bazelrc` for a full list of available namespaced parameters.

## Usage

To generate `third_party/grlib/grlib.BUILD` files for the libraries, run:

```bash
bazel run //third_party/grlib/scripts:gen_build_files
```

To build all targets:

```bash
bazel build //...
```
