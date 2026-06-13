<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="ElaborateProvider"></a>

## ElaborateProvider

<pre>
load("@rules_nvc//internal:providers.bzl", "ElaborateProvider")

ElaborateProvider(<a href="#ElaborateProvider-entity">entity</a>)
</pre>

Provides information about an elaborated VHDL entity.

**FIELDS**

| Name  | Description |
| :------------- | :------------- |
| <a id="ElaborateProvider-entity"></a>entity |  string: The name of the elaborated entity.    |


<a id="NVCInfo"></a>

## NVCInfo

<pre>
load("@rules_nvc//internal:providers.bzl", "NVCInfo")

NVCInfo(<a href="#NVCInfo-analyzer">analyzer</a>, <a href="#NVCInfo-artifacts_dir">artifacts_dir</a>, <a href="#NVCInfo-deps">deps</a>)
</pre>

Information on how to run NVC for VHDL analysis, elaboration and simulation.

**FIELDS**

| Name  | Description |
| :------------- | :------------- |
| <a id="NVCInfo-analyzer"></a>analyzer |  The NVC analyzer executable file.    |
| <a id="NVCInfo-artifacts_dir"></a>artifacts_dir |  The directory containing NVC standard libraries and artifacts.    |
| <a id="NVCInfo-deps"></a>deps |  List[File]: Additional toolchain dependencies (e.g. shared libraries).    |


<a id="VHDLLibraryProvider"></a>

## VHDLLibraryProvider

<pre>
load("@rules_nvc//internal:providers.bzl", "VHDLLibraryProvider")

VHDLLibraryProvider(<a href="#VHDLLibraryProvider-libraries">libraries</a>, <a href="#VHDLLibraryProvider-entities">entities</a>, <a href="#VHDLLibraryProvider-library_name">library_name</a>, <a href="#VHDLLibraryProvider-library_dir">library_dir</a>, <a href="#VHDLLibraryProvider-includes">includes</a>, <a href="#VHDLLibraryProvider-hdrs">hdrs</a>, <a href="#VHDLLibraryProvider-vpi_plugins">vpi_plugins</a>)
</pre>

Contains the information about the binary files in a compiled VHDL library.

**FIELDS**

| Name  | Description |
| :------------- | :------------- |
| <a id="VHDLLibraryProvider-libraries"></a>libraries |  List[(string, File)]: A mapping from a library name to its directory location. Contains both this library and its dependencies, ensuring no duplicate keys.    |
| <a id="VHDLLibraryProvider-entities"></a>entities |  List[string]: The entities emphasized in this library.    |
| <a id="VHDLLibraryProvider-library_name"></a>library_name |  string: The name of the library (e.g., `ieee`, `work`).    |
| <a id="VHDLLibraryProvider-library_dir"></a>library_dir |  File: The container directory where the library is located. NVC will not create a library in an existing directory, so a container directory is used. The actual library directory is `$library_dir/$library_name`.    |
| <a id="VHDLLibraryProvider-includes"></a>includes |  List[string]: List of directories to include for Verilog.    |
| <a id="VHDLLibraryProvider-hdrs"></a>hdrs |  depset[File]: List of include files for Verilog.    |
| <a id="VHDLLibraryProvider-vpi_plugins"></a>vpi_plugins |  depset[File]: List of VPI plugins required for simulation.    |


