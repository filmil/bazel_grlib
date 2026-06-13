def _patched_nvc_impl(rctx):
    # 1. Download official NVC r1.21.0 release source
    rctx.report_progress("Downloading NVC r1.21.0 source")
    rctx.download_and_extract(
        url = "https://github.com/nickg/nvc/archive/refs/tags/r1.21.0.tar.gz",
        stripPrefix = "nvc-r1.21.0",
    )

    # 2. Apply our custom GRLIB patches
    rctx.report_progress("Applying GRLIB compiler patches")
    rctx.patch(rctx.attr.patch, strip = 1)

    # 3. Build NVC from source using the host compiler
    rctx.report_progress("Configuring NVC")
    res = rctx.execute(["autoreconf", "--force", "--install", "-I", "m4"])
    if res.return_code != 0:
        fail("autoreconf failed: " + res.stderr)

    # Create the out-of-tree build directory
    rctx.execute(["mkdir", "-p", "build"])

    # Run configure out-of-tree
    res = rctx.execute(
        ["../configure", "--prefix=" + str(rctx.path("usr"))],
        working_directory = "build"
    )
    if res.return_code != 0:
        fail("configure failed: " + res.stderr)

    rctx.report_progress("Compiling NVC")
    res = rctx.execute(
        ["make", "-j4"],
        working_directory = "build"
    )
    if res.return_code != 0:
        fail("make failed: " + res.stderr)

    rctx.report_progress("Installing NVC")
    res = rctx.execute(
        ["make", "install"],
        working_directory = "build"
    )
    if res.return_code != 0:
        fail("make install failed: " + res.stderr)

    # Create the expected standard usr/lib/x86_64-linux-gnu symlink so existing toolchain configurations work perfectly
    rctx.execute(["mkdir", "-p", "usr/lib/x86_64-linux-gnu"])
    rctx.symlink("usr/lib/nvc", "usr/lib/x86_64-linux-gnu/nvc")

    # Create the BUILD.bazel file
    rctx.file(
        "BUILD.bazel",
        """package(default_visibility = ["//visibility:public"])
exports_files(glob(['**/*']))
filegroup(
    name = "all_files",
    srcs = glob(['**/*']),
    visibility = ["//visibility:public"]
)"""
    )

patched_nvc = repository_rule(
    implementation = _patched_nvc_impl,
    attrs = {
        "patch": attr.label(mandatory = True),
    },
)
