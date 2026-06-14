# GRLIB Bazel Repository

![Build Status](https://github.com/filmil/bazel_grlib/actions/workflows/build.yml/badge.svg)
![Publish to my Bazel registry](https://github.com/filmil/bazel-registry/main/badge.svg)
![Publish on Bazel Central Registry](https://bcr.bazel.build/badge.svg)
![Tag and Release](https://github.com/filmil/bazel_grlib/actions/workflows/tag-and-release.yml/badge.svg)

**Note:** The repository name is `bazel_grlib`, while the Bazel module name is `grlib`.

This repository contains a Bazel-based build system for GRLIB VHDL designs.

## Architecture

The build system is designed to be tool-agnostic. All VHDL source files are organized into `filegroup` targets within the `@grlib_srcs` repository. Tool-specific rules (e.g., for the `nvc` compiler) are isolated in the `tool.nvc` subdirectory.

## Configuration

This repository uses `rules_kconfig` to manage hardware configuration. Every GRLIB component's configuration is namespaced to avoid collisions.

### Namespacing

Symbols are prefixed with their GRLIB path. For example, `IU_NWINDOWS` defined in `lib/gaisler/leon3/leon3.in` becomes:

`@grlib_config//:CONFIG_LIB_GAISLER_LEON3_LEON3_IU_NWINDOWS`

### Promoting a Design Configuration

To make a specific design's configuration the "active" one (mapping namespaced symbols like `DESIGNS_LEON3MP_IU_NWINDOWS` back to the generic `CFG_IU_NWINDOWS` expected by VHDL), set the `CONFIG_ACTIVE_DESIGN_PREFIX`:

```bash
# Example: Activate the LEON3MP design configuration
bazel build --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=DESIGNS_LEON3MP //...
```

### Specifying the VHDL Version (Standard)

By default, the project uses VHDL 1993. You can specify a different VHDL version or standard (e.g., 2008 or 2019) using the `--@grlib//:vhdl_standard` flag. This version flag controls the standard used for compilation:

```bash
bazel test //tool.nvc:noelv_simulation --@grlib//:vhdl_standard=2019
```

Available versions are: `1987`, `1993`, `2002`, `2008`, `2019`.

## Usage

To generate the build files (including tool-agnostic filegroups and NVC rules), run:

```bash
bazel run //third_party/grlib/scripts:gen_build_files
```

To run the NOEL-V simulation test (using NVC):

```bash
bazel test //tool.nvc:noelv_simulation --@grlib_config//:CONFIG_ACTIVE_DESIGN_PREFIX=LIB_GAISLER_NOELV_NOELV
```

### Reference

See `docs/worked_example.md` for a step-by-step guide on customizing your build, and `docs/example.bazelrc` for a full list of available namespaced parameters.
