<!-- LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4 -->

# rules_nvc

[![Build](https://github.com/filmil/bazel_rules_nvc/actions/workflows/build.yml/badge.svg)](https://github.com/filmil/bazel_rules_nvc/actions/workflows/build.yml)
[![Publish on Bazel Central Registry](https://github.com/filmil/bazel_rules_nvc/actions/workflows/publish-bcr.yml/badge.svg)](https://github.com/filmil/bazel_rules_nvc/actions/workflows/publish-bcr.yml)
[![Publish to my Bazel registry](https://github.com/filmil/bazel_rules_nvc/actions/workflows/publish.yml/badge.svg)](https://github.com/filmil/bazel_rules_nvc/actions/workflows/publish.yml)
[![Tag and Release](https://github.com/filmil/bazel_rules_nvc/actions/workflows/tag-and-release.yml/badge.svg)](https://github.com/filmil/bazel_rules_nvc/actions/workflows/tag-and-release.yml)

This repository provides [Bazel][baz] build rules for the [NVC compiler and
simulator][nvc]. NVC is a VHDL compiler and simulator that aims for VHDL-2019
compliance. These rules allow you to integrate NVC into your Bazel-based VHDL
projects, enabling you to build, simulate, and test your VHDL code within the
Bazel ecosystem.

[baz]: https://bazel.build
[nvc]: https://github.com/nickg/nvc

## Documentation

An index of architectural reports and feature documentation can be found in the [Documentation Index](doc/README.md).

Key documentation includes:
- [Co-Simulation Usage Example](doc/cosim_usage_example.md)
- [Simulation vs. Synthesis Flow](doc/simulation_vs_synthesis.md)

The following table provides an overview of the documentation for the Starlark
rules and macros provided by this repository.

| File | Description |
| :--- | :--- |
| [macros.md](macros.md) | Top-level macros |
| [nvc/rules.md](nvc/rules.md) | User-facing NVC rules |
| [build/nvc/rules.md](build/nvc/rules.md) | NVC toolchain rules |
| [build/nvc/repositories.md](build/nvc/repositories.md) | NVC repository rules |
| [internal/vhdl_library.md](internal/vhdl_library.md) | VHDL library compilation |
| [internal/vhdl_elaborate.md](internal/vhdl_elaborate.md) | VHDL elaboration |
| [internal/vhdl_run.md](internal/vhdl_run.md) | VHDL simulation execution |
| [internal/vhdl_test.md](internal/vhdl_test.md) | VHDL test execution |
| [internal/produce_waveform.md](internal/produce_waveform.md) | Waveform production |
| [internal/wave_view.md](internal/macros.md) | Waveform viewing |
| [internal/extract_file.md](internal/extract_file.md) | File extraction utilities |
| [internal/prebuilt_library.md](internal/prebuilt_library.md) | Prebuilt library support |
| [internal/verilog_library.md](internal/verilog_library.md) | Verilog library support |
| [internal/toolchain.md](internal/toolchain.md) | Toolchain implementation |
| [internal/providers.md](internal/providers.md) | Starlark providers |
| [internal/utils.md](internal/utils.md) | Internal utilities |
| [third_party/repositories.md](third_party/repositories.md) | Third-party repository rules |

## Hermeticity

* Hermeticity in this repository is best effort.
* `nvc` is built from source. This may last a *long* time.
* At the moment I believe that the remaining non-hermetic dependency
  is `llvm-config`.
