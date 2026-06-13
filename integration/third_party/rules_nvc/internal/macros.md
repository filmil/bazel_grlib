<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="wave_view"></a>

## wave_view

<pre>
load("@rules_nvc//internal:macros.bzl", "wave_view")

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


