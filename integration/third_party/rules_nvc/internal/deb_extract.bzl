def _deb_extract_impl(rctx):
    rctx.download(
        url = rctx.attr.urls,
        integrity = rctx.attr.integrity,
        output = "package.deb",
    )
    # Use ar to extract the deb. We assume ar is available on the host.
    res = rctx.execute(["ar", "x", "package.deb"])
    if res.return_code != 0:
        fail("Failed to extract .deb: " + res.stderr)

    # Identify the data archive. It could be data.tar.xz, data.tar.gz, or data.tar.zst
    data_archive = None
    for f in ["data.tar.zst", "data.tar.xz", "data.tar.gz"]:
        if rctx.path(f).exists:
            data_archive = f
            break

    if not data_archive:
        fail("Could not find data archive in .deb")

    rctx.extract(data_archive)

    if rctx.attr.build_file:
        rctx.symlink(rctx.attr.build_file, "BUILD.bazel")
    elif rctx.attr.build_file_content:
        rctx.file("BUILD.bazel", rctx.attr.build_file_content)
    else:
        # Default BUILD file that exports everything
        rctx.file("BUILD.bazel", "exports_files(glob(['**/*']))\nfilegroup(name = 'all_files', srcs = glob(['**/*']), visibility = ['//visibility:public'])")

deb_extract = repository_rule(
    implementation = _deb_extract_impl,
    attrs = {
        "urls": attr.string_list(mandatory = True),
        "integrity": attr.string(mandatory = True),
        "build_file": attr.label(),
        "build_file_content": attr.string(),
    },
)
