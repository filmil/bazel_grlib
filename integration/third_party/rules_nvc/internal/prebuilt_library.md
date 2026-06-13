<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="prebuilt_library"></a>

## prebuilt_library

<pre>
load("@rules_nvc//internal:prebuilt_library.bzl", "prebuilt_library")

prebuilt_library(<a href="#prebuilt_library-name">name</a>, <a href="#prebuilt_library-deps">deps</a>, <a href="#prebuilt_library-srcs">srcs</a>, <a href="#prebuilt_library-container_dir">container_dir</a>, <a href="#prebuilt_library-entities">entities</a>, <a href="#prebuilt_library-library_name">library_name</a>, <a href="#prebuilt_library-library_remapped_name">library_remapped_name</a>,
                 <a href="#prebuilt_library-standard">standard</a>)
</pre>

Defines a prebuilt VHDL library, skipping analysis.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="prebuilt_library-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="prebuilt_library-deps"></a>deps |  The dependency libraries, if any   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="prebuilt_library-srcs"></a>srcs |  The list of source targets that comprise this prebuilt library.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="prebuilt_library-container_dir"></a>container_dir |  The library directory, used to short-circuit directory detection   | <a href="https://bazel.build/concepts/labels">Label</a> | required |  |
| <a id="prebuilt_library-entities"></a>entities |  The list of entities emphasized in thsi library, for elaboration purposes   | List of strings | optional |  `[]`  |
| <a id="prebuilt_library-library_name"></a>library_name |  The official library name, in case target name is not appropriate. Target name is used if not specified.   | String | optional |  `""`  |
| <a id="prebuilt_library-library_remapped_name"></a>library_remapped_name |  Usually `unisim.08` for `unisim` when used with `standard=2008`   | String | optional |  `""`  |
| <a id="prebuilt_library-standard"></a>standard |  The HDL language standard to use. For VHDL it is the standard version, such as 1993, or 2008, or 2019   | String | optional |  `""`  |


