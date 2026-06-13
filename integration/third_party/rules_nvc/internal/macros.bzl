load("@rules_shell//shell:sh_binary.bzl", "sh_binary")


def wave_view(name, vhdl_run, args=[], deps=[], viewer="gtkwave", testonly=None, save_file=None):
    """
    Generates a sh_binary viewer.

    # Args

    - name: the target name.
    - vhdl_run: the target name for the `vhdl_run` target to
      use the output from.
    - args: any additional arguments to add to invoke the viewer.
    - viewer: the viewer to invoke. The viewer must be compatible
      with the file format to view, and must be installed on the
      host.

    """
    _args = [
          "--viewer-binary={}".format(viewer),
          "--wave-file=$(location {})".format(vhdl_run),
          "--",
    ] + args
    _data = [
          "@gotopt2//:bin",
          vhdl_run,
        ] + deps
    if save_file:
        _args += ["--save=$(location {})".format(save_file)]
        _data += [save_file]
    sh_binary(
        testonly = testonly,
        name = name,
        srcs = [
            # This is necessary because run_wave_view below is multi-file, and
            # multi-file can not be in `sh_binary` srcs.
            # Using `Label` here so that the string "file_shunt.bash" refers
            # to the file where macro is defined, not the file in the dir where
            # the macro is invoked.
            Label("file_shunt.bash"),
        ],
        deps = [
            "@bazel_tools//tools/bash/runfiles",
            Label("//build/nvc:run_wave_view"),
            "@fshlib//:log",
        ],
        args = _args,
        data = _data,
        visibility = ["//visibility:public"],
    )
