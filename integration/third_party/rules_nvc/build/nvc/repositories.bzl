load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(
    "@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")
load("@bazel_skylib//lib:modules.bzl", "modules")


def nvc_repositories():
    """
    Declares the external repositories required by the NVC toolchain.

    This function sets up the `http_archive` for the NVC compiler source.
    """
    maybe(
        repo_rule = http_archive,
        name = "nvc",
        build_file = Label("//third_party/nvc:nvc.BUILD.bazel"),
        patch_args = ["-p1"],
        sha256 = "192fe81768d76d90ea005dcde1ad997ec5220a5b84103c763f39758f12cbb4a3",
        strip_prefix = "nvc-1.11.0",
        urls = [
            "https://github.com/nickg/nvc/releases/download/r1.11.0/nvc-1.11.0.tar.gz",
        ],
    )


repositories_extension = modules.as_extension(nvc_repositories)

