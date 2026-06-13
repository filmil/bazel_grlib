#!/usr/bin/env bash
# GENERATED FILE DO NOT EDIT.
#
# Library name:  {{LIBRARY_NAME}}
# Entity:        {{ENTITY}}
# VHDL Standard: {{VHDL_STANDARD}}
set -eo pipefail
export TMPDIR="$TEST_TMPDIR"

readonly wave_file="${TEST_UNDECLARED_OUTPUTS_DIR}/{{WAVE_FILE}}"

readonly dir_in_path="{{LIB_DIR_IN_PATH}}"
mkdir -p "${dir_in_path}"

readonly dir_out_path="${TEST_TMPDIR}/{{LIB_DIR_OUT_PATH}}"
mkdir -p "${dir_out_path}/{{LIBRARY_NAME}}"

{{EXECUTABLE}} -cmd=-r \
    --vhdl-standard={{VHDL_STANDARD}} \
    --nvc-binary-path={{ANALYZER}} \
    --library-name={{LIBRARY_NAME}} \
    --library-paths="{{LIBRARY_PATHS}}" \
    --stdlib-dir={{STDLIB_DIR}} \
    --entity={{ENTITY}} \
    --library-dir-in-path="${dir_in_path}" \
    --library-dir-out-path="${dir_in_path}" \
    -- \
    {{VPI_FLAGS}} \
    {{EXTRA_ARGS}} \
    --wave="${wave_file}" \
    $@

