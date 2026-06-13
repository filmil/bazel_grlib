<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="get_nvc_deps"></a>

## get_nvc_deps

<pre>
load("@rules_nvc//internal:utils.bzl", "get_nvc_deps")

get_nvc_deps(<a href="#get_nvc_deps-nvc_info">nvc_info</a>)
</pre>

Returns the list of dependencies from nvc_info.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="get_nvc_deps-nvc_info"></a>nvc_info |  <p align="center"> - </p>   |  none |


<a id="get_nvc_ld_library_path"></a>

## get_nvc_ld_library_path

<pre>
load("@rules_nvc//internal:utils.bzl", "get_nvc_ld_library_path")

get_nvc_ld_library_path(<a href="#get_nvc_ld_library_path-nvc_info">nvc_info</a>, <a href="#get_nvc_ld_library_path-base_dir">base_dir</a>, <a href="#get_nvc_ld_library_path-default_env">default_env</a>)
</pre>

Constructs the LD_LIBRARY_PATH environment variable.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="get_nvc_ld_library_path-nvc_info"></a>nvc_info |  <p align="center"> - </p>   |  none |
| <a id="get_nvc_ld_library_path-base_dir"></a>base_dir |  <p align="center"> - </p>   |  none |
| <a id="get_nvc_ld_library_path-default_env"></a>default_env |  <p align="center"> - </p>   |  none |


<a id="get_single_file_from"></a>

## get_single_file_from

<pre>
load("@rules_nvc//internal:utils.bzl", "get_single_file_from")

get_single_file_from(<a href="#get_single_file_from-target">target</a>)
</pre>

Retrieves the single file associated with a target.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="get_single_file_from-target"></a>target |  The target from which to extract the file.   |  none |

**RETURNS**

The single `File` object from the target.


