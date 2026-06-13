<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="vhdl_library"></a>

## vhdl_library

<pre>
load("@rules_nvc//internal:vhdl_library.bzl", "vhdl_library")

vhdl_library(<a href="#vhdl_library-name">name</a>, <a href="#vhdl_library-deps">deps</a>, <a href="#vhdl_library-srcs">srcs</a>, <a href="#vhdl_library-entities">entities</a>, <a href="#vhdl_library-library_name">library_name</a>, <a href="#vhdl_library-standard">standard</a>, <a href="#vhdl_library-vpi_plugins">vpi_plugins</a>)
</pre>

Compiles VHDL source files into a library using NVC.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="vhdl_library-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="vhdl_library-deps"></a>deps |  A list of other `vhdl_library` targets that this library depends on.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="vhdl_library-srcs"></a>srcs |  A list of VHDL source files.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="vhdl_library-entities"></a>entities |  A list of VHDL entities provided by this library.   | List of strings | optional |  `[]`  |
| <a id="vhdl_library-library_name"></a>library_name |  If the target name is not appropriate as a library name, provide one here   | String | optional |  `""`  |
| <a id="vhdl_library-standard"></a>standard |  The VHDL standard to use for compilation (e.g., '2008', '2019'). Defaults to '2019'.   | String | optional |  `"2019"`  |
| <a id="vhdl_library-vpi_plugins"></a>vpi_plugins |  List of VPI plugins required for simulation.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |


