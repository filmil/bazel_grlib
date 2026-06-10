# GRLIB Bazel Repository

![Build Status](https://github.com/filmil/bazel_grlib/actions/workflows/build.yml/badge.svg)

**Note:** The repository name is `bazel_grlib`, while the Bazel module name is `grlib`.

This repository contains a Bazel-based build system for GRLIB VHDL designs.

## Configuration

This repository uses `rules_kconfig` to manage hardware configuration. You can define your hardware parameters in the root `Kconfig` file.

To set configuration values, you can use Bazel command-line flags. Each Kconfig symbol is mapped to a Bazel build setting in the `@grlib_config` repository.

### Setting values via command line

For a boolean option `CONFIG_GRLIB_DUMMY`:

```bash
# To enable (True)
bazel build --@grlib_config//:CONFIG_GRLIB_DUMMY=True //...

# To disable (False)
bazel build --@grlib_config//:CONFIG_GRLIB_DUMMY=False //...
```

For integer or string options:

```bash
bazel build --@grlib_config//:CONFIG_MY_INT=42 //...
bazel build --@grlib_config//:CONFIG_MY_STRING="hello" //...
```

### Setting values via .bazelrc

You can also add these flags to your `user.bazelrc` or `.bazelrc`:

```text
build --@grlib_config//:CONFIG_GRLIB_DUMMY=True
```

## Usage

To generate `third_party/grlib/grlib.BUILD` files for the libraries, run:

```bash
bazel run //third_party/grlib/scripts:gen_build_files
```

To build all targets:

```bash
bazel build //...
```
