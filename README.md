# GRLIB Bazel Repository

![Build Status](https://github.com/filmil/bazel_grlib/actions/workflows/build.yml/badge.svg)

**Note:** The repository name is `bazel_grlib`, while the Bazel module name is `grlib`.

This repository contains a Bazel-based build system for GRLIB VHDL designs.

## Usage

To generate `third_party/grlib/grlib.BUILD` files for the libraries, run:

```bash
bazel run //third_party/grlib/scripts:gen_build_files
```

To build all targets:

```bash
bazel build //...
```
