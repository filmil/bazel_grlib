<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="verilog_library"></a>

## verilog_library

<pre>
load("@rules_nvc//internal:verilog_library.bzl", "verilog_library")

verilog_library(<a href="#verilog_library-name">name</a>, <a href="#verilog_library-deps">deps</a>, <a href="#verilog_library-srcs">srcs</a>, <a href="#verilog_library-hdrs">hdrs</a>, <a href="#verilog_library-entities">entities</a>, <a href="#verilog_library-includes">includes</a>, <a href="#verilog_library-library_name">library_name</a>, <a href="#verilog_library-standard">standard</a>)
</pre>

Compiles Verilog source files into a library using NVC.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="verilog_library-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="verilog_library-deps"></a>deps |  List of dependency libraries.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="verilog_library-srcs"></a>srcs |  List of Verilog source files.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="verilog_library-hdrs"></a>hdrs |  List of Verilog header files.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="verilog_library-entities"></a>entities |  List of entities provided by this library.   | List of strings | optional |  `[]`  |
| <a id="verilog_library-includes"></a>includes |  list of verilog include directories   | List of strings | optional |  `[]`  |
| <a id="verilog_library-library_name"></a>library_name |  If the target name is not appropriate as a library name, provide one here   | String | optional |  `""`  |
| <a id="verilog_library-standard"></a>standard |  The VHDL standard to use for compilation.   | String | optional |  `"2019"`  |


