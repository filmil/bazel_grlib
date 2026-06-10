# GRLIB Bazel Repository

![Build Status](https://github.com/filmil/grlib/actions/workflows/build.yml/badge.svg)

This repository contains a Bazel-based build system for GRLIB VHDL designs.

## Usage

To generate `BUILD.bazel` files for the libraries, run:

```bash
bazel run //scripts:gen_build_files
```

To build all targets:

```bash
bazel build //...
```
