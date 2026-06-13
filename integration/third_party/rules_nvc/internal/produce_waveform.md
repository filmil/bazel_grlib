<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="produce_waveform"></a>

## produce_waveform

<pre>
load("@rules_nvc//internal:produce_waveform.bzl", "produce_waveform")

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


