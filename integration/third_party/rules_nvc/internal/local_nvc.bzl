def _local_nvc_impl(rctx):
    # Create usr/bin directory
    rctx.execute(["mkdir", "-p", "usr/bin"])
    # Symlink our compiled nvc binary
    rctx.symlink(rctx.attr.nvc_binary_path, "usr/bin/nvc")

    # Create usr/lib/x86_64-linux-gnu directory
    rctx.execute(["mkdir", "-p", "usr/lib/x86_64-linux-gnu"])
    # Symlink our compiled libraries directory
    rctx.symlink(rctx.attr.nvc_lib_dir, "usr/lib/x86_64-linux-gnu/nvc")

    # Default BUILD file that exports everything
    rctx.file(
        "BUILD.bazel",
        "exports_files(glob(['**/*']))\nfilegroup(name = 'all_files', srcs = glob(['**/*']), visibility = ['//visibility:public'])"
    )

local_nvc = repository_rule(
    implementation = _local_nvc_impl,
    attrs = {
        "nvc_binary_path": attr.string(mandatory = True),
        "nvc_lib_dir": attr.string(mandatory = True),
    },
)
