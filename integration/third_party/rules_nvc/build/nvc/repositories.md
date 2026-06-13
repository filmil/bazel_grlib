<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="nvc_repositories"></a>

## nvc_repositories

<pre>
load("@rules_nvc//build/nvc:repositories.bzl", "nvc_repositories")

nvc_repositories()
</pre>

Declares the external repositories required by the NVC toolchain.

This function sets up the `http_archive` for the NVC compiler source.



<a id="repositories_extension"></a>

## repositories_extension

<pre>
repositories_extension = use_extension("@rules_nvc//build/nvc:repositories.bzl", "repositories_extension")
</pre>



