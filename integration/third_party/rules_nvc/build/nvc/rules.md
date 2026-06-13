<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="extract_file"></a>

## extract_file

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "extract_file")

extract_file(<a href="#extract_file-name">name</a>, <a href="#extract_file-src">src</a>, <a href="#extract_file-filter">filter</a>)
</pre>

Extracts specific files from a source target based on a string filter.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="extract_file-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="extract_file-src"></a>src |  The source target to extract files from.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |
| <a id="extract_file-filter"></a>filter |  A string used to filter the files in the source target.   | String | optional |  `""`  |


<a id="nvc_toolchain"></a>

## nvc_toolchain

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "nvc_toolchain")

nvc_toolchain(<a href="#nvc_toolchain-name">name</a>, <a href="#nvc_toolchain-deps">deps</a>, <a href="#nvc_toolchain-analyzer">analyzer</a>, <a href="#nvc_toolchain-artifacts_dir">artifacts_dir</a>)
</pre>

Defines the NVC toolchain, linking to the NVC analyzer and standard library.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="nvc_toolchain-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="nvc_toolchain-deps"></a>deps |  Additional toolchain dependencies.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="nvc_toolchain-analyzer"></a>analyzer |  The NVC executable wrapper script.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |
| <a id="nvc_toolchain-artifacts_dir"></a>artifacts_dir |  The directory containing NVC standard libraries and artifacts.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |


<a id="prebuilt_library"></a>

## prebuilt_library

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "prebuilt_library")

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


<a id="produce_waveform"></a>

## produce_waveform

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "produce_waveform")

produce_waveform(<a href="#produce_waveform-name">name</a>, <a href="#produce_waveform-data">data</a>, <a href="#produce_waveform-args">args</a>, <a href="#produce_waveform-simulation">simulation</a>, <a href="#produce_waveform-use_fst">use_fst</a>)
</pre>

Produces a waveform file (VCD) from a VHDL simulation run.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="produce_waveform-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="produce_waveform-data"></a>data |  Data files required for the simulation.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="produce_waveform-args"></a>args |  Additional command-line arguments to pass to the simulation.   | List of strings | optional |  `[]`  |
| <a id="produce_waveform-simulation"></a>simulation |  The simulation target (`vhdl_run`) to execute.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |
| <a id="produce_waveform-use_fst"></a>use_fst |  A boolean indicating whether to expect an FST file instead of VCD file. Defaults to `False`.   | Boolean | optional |  `False`  |


<a id="verilog_library"></a>

## verilog_library

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "verilog_library")

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


<a id="vhdl_elaborate"></a>

## vhdl_elaborate

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "vhdl_elaborate")

vhdl_elaborate(<a href="#vhdl_elaborate-name">name</a>, <a href="#vhdl_elaborate-library">library</a>, <a href="#vhdl_elaborate-standard">standard</a>)
</pre>

Elaborates a VHDL design using NVC.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="vhdl_elaborate-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="vhdl_elaborate-library"></a>library |  The `vhdl_library` target to elaborate.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |
| <a id="vhdl_elaborate-standard"></a>standard |  The VHDL standard to use for elaboration (e.g., '2019').   | String | optional |  `"2019"`  |


<a id="vhdl_library"></a>

## vhdl_library

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "vhdl_library")

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


<a id="vhdl_run"></a>

## vhdl_run

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "vhdl_run")

vhdl_run(<a href="#vhdl_run-name">name</a>, <a href="#vhdl_run-deps">deps</a>, <a href="#vhdl_run-args">args</a>, <a href="#vhdl_run-entity">entity</a>, <a href="#vhdl_run-standard">standard</a>, <a href="#vhdl_run-use_fst">use_fst</a>, <a href="#vhdl_run-use_vcd">use_vcd</a>)
</pre>

Simulates an elaborated VHDL design using NVC.

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="vhdl_run-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="vhdl_run-deps"></a>deps |  A list of other `vhdl_library` targets that this simulation depends on.   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="vhdl_run-args"></a>args |  A list of added command line args to use   | List of strings | optional |  `[]`  |
| <a id="vhdl_run-entity"></a>entity |  The elaborated VHDL entity to simulate. This should be a `vhdl_elaborate` target.   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `None`  |
| <a id="vhdl_run-standard"></a>standard |  The VHDL standard to use for simulation. Defaults to '2019'.   | String | optional |  `"2019"`  |
| <a id="vhdl_run-use_fst"></a>use_fst |  A boolean indicating whether to generate a FST file for waveform viewing. Defaults to `False`. Takes precedence over `use_vcd`.   | Boolean | optional |  `False`  |
| <a id="vhdl_run-use_vcd"></a>use_vcd |  A boolean indicating whether to generate a VCD (Value Change Dump) file for waveform viewing. Defaults to `True`.   | Boolean | optional |  `True`  |


<a id="vhdl_test"></a>

## vhdl_test

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "vhdl_test")

vhdl_test(<a href="#vhdl_test-name">name</a>, <a href="#vhdl_test-srcs">srcs</a>, <a href="#vhdl_test-deps">deps</a>, <a href="#vhdl_test-standard">standard</a>, <a href="#vhdl_test-args">args</a>, <a href="#vhdl_test-entity">entity</a>, <a href="#vhdl_test-entities">entities</a>)
</pre>

Defines a VHDL test.

This macro combines `vhdl_library`, `vhdl_elaborate`, and internal test
execution steps into a single logical target.


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="vhdl_test-name"></a>name |  The name of the base test target.   |  none |
| <a id="vhdl_test-srcs"></a>srcs |  A list of VHDL source files (`.vhdl` or `.vhd`).   |  none |
| <a id="vhdl_test-deps"></a>deps |  A list of `vhdl_library` targets that this test depends on.   |  none |
| <a id="vhdl_test-standard"></a>standard |  The VHDL standard to use (e.g., "2008", "2019"). Defaults to "2019".   |  `"2019"` |
| <a id="vhdl_test-args"></a>args |  A list of additional command-line arguments to pass to the NVC simulator.   |  `[]` |
| <a id="vhdl_test-entity"></a>entity |  A single entity to test.   |  `None` |
| <a id="vhdl_test-entities"></a>entities |  A list of entities to test. If both `entity` and `entities` are provided, all are tested.   |  `[]` |


<a id="wave_view"></a>

## wave_view

<pre>
load("@rules_nvc//build/nvc:rules.bzl", "wave_view")

wave_view(<a href="#wave_view-name">name</a>, <a href="#wave_view-vhdl_run">vhdl_run</a>, <a href="#wave_view-args">args</a>, <a href="#wave_view-deps">deps</a>, <a href="#wave_view-viewer">viewer</a>, <a href="#wave_view-testonly">testonly</a>, <a href="#wave_view-save_file">save_file</a>)
</pre>

Generates a sh_binary viewer.

# Args

- name: the target name.
- vhdl_run: the target name for the `vhdl_run` target to
  use the output from.
- args: any additional arguments to add to invoke the viewer.
- viewer: the viewer to invoke. The viewer must be compatible
  with the file format to view, and must be installed on the
  host.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="wave_view-name"></a>name |  <p align="center"> - </p>   |  none |
| <a id="wave_view-vhdl_run"></a>vhdl_run |  <p align="center"> - </p>   |  none |
| <a id="wave_view-args"></a>args |  <p align="center"> - </p>   |  `[]` |
| <a id="wave_view-deps"></a>deps |  <p align="center"> - </p>   |  `[]` |
| <a id="wave_view-viewer"></a>viewer |  <p align="center"> - </p>   |  `"gtkwave"` |
| <a id="wave_view-testonly"></a>testonly |  <p align="center"> - </p>   |  `None` |
| <a id="wave_view-save_file"></a>save_file |  <p align="center"> - </p>   |  `None` |


