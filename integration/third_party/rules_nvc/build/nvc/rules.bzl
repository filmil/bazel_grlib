load("//internal:utils.bzl", "get_single_file_from")
load("//internal:providers.bzl", "NVCInfo", "VHDLLibraryProvider", "ElaborateProvider")
load("//internal:toolchain.bzl",
     _NVC_TOOLCHAIN_TYPE = "NVC_TOOLCHAIN_TYPE",
     _NVC_WRAPPER = "NVC_WRAPPER",
    _VHDL_STANDARD_DEFAULT = "VHDL_STANDARD_DEFAULT",
    _nvc_toolchain = "nvc_toolchain")


load("//internal:extract_file.bzl", _extract_file = "extract_file")
load("//internal:macros.bzl", _wave_view = "wave_view")
load("//internal:prebuilt_library.bzl", _prebuilt_library = "prebuilt_library")
load("//internal:produce_waveform.bzl", _produce_waveform = "produce_waveform")
load("//internal:vhdl_elaborate.bzl", _vhdl_elaborate = "vhdl_elaborate")
load("//internal:vhdl_library.bzl", _vhdl_library = "vhdl_library")
load("//internal:vhdl_run.bzl", _vhdl_run = "vhdl_run")
load("//internal:vhdl_test.bzl", _vhdl_test = "vhdl_test")

load("//internal:verilog_library.bzl", _verilog_library = "verilog_library")


# The main API
nvc_toolchain = _nvc_toolchain
vhdl_library = _vhdl_library
vhdl_elaborate = _vhdl_elaborate
vhdl_run = _vhdl_run
produce_waveform = _produce_waveform
extract_file = _extract_file
prebuilt_library = _prebuilt_library
vhdl_test = _vhdl_test

verilog_library = _verilog_library


# Macros
wave_view = _wave_view

