#! /bin/bash

# Copyright (C) 2020 Google Inc.
#
# This file has been licensed under Apache 2.0 license.  Please see the LICENSE
# file at the root of the repository.

set -e

# This magic was copied from runfiles by consulting:
#   https://stackoverflow.com/questions/53472993/how-do-i-make-a-bazel-sh-binary-target-depend-on-other-binary-targets

# --- begin runfiles.bash initialization ---
# Copy-pasted from Bazel's Bash runfiles library (tools/bash/runfiles/runfiles.bash).
set -eo pipefail
if [[ ! -d "${RUNFILES_DIR:-/dev/null}" && ! -f "${RUNFILES_MANIFEST_FILE:-/dev/null}" ]]; then
  if [[ -f "$0.runfiles_manifest" ]]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles_manifest"
  elif [[ -f "$0.runfiles/MANIFEST" ]]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles/MANIFEST"
  elif [[ -f "$0.runfiles/bazel_tools/tools/bash/runfiles/runfiles.bash" ]]; then
    export RUNFILES_DIR="$0.runfiles"
  fi
fi
if [[ -f "${RUNFILES_DIR:-/dev/null}/bazel_tools/tools/bash/runfiles/runfiles.bash" ]]; then
  source "${RUNFILES_DIR}/bazel_tools/tools/bash/runfiles/runfiles.bash"
elif [[ -f "${RUNFILES_MANIFEST_FILE:-/dev/null}" ]]; then
  source "$(grep -m1 "^bazel_tools/tools/bash/runfiles/runfiles.bash " \
            "$RUNFILES_MANIFEST_FILE" | cut -d ' ' -f 2-)"
else
  echo >&2 "ERROR: cannot find @bazel_tools//tools/bash/runfiles:runfiles.bash"
  exit 1
fi
# --- end runfiles.bash initialization ---
set -u


_gotopt2_runfiles_paths=(
  "rules_multitool++multitool+multitool/tools/gotopt2/gotopt2"
  "rules_multitool~~multitool~multitool/tools/gotopt2/gotopt2"
  "gotopt2/bin/gotopt2"
)
_gotopt_binary=""
for binary_candidate in ${_gotopt2_runfiles_paths[@]}; do
  _gotopt_binary="$(rlocation ${binary_candidate})"
  if [[ -x "${_gotopt_binary}" ]]; then
    break
  fi
done
# Exit quickly if the binary isn't found. This may happen if the binary location
# moves internally in bazel.
if [[ ! -x "$(command -v ${_gotopt_binary})" ]]; then
  echo "gotopt2 binary not found"
  exit 240
fi


GOTOPT2_OUTPUT=$(${_gotopt_binary} "${@}" <<EOF
flags:
- name: "viewer-binary"
  type: string
  default: "gtkwave"
  help: "the binary to start to view waveforms"
- name: "wave-file"
  type: string
  help: "the wave file to display"
EOF
)
if [[ "$?" == "11" ]]; then
  # When --help option is used, gotopt2 exits with code 11.
  exit 1
fi

# Evaluate the output of the call to gotopt2, shell vars assignment is here.
eval "${GOTOPT2_OUTPUT}"

# Exit if required parameters are not present.
if [[ "${gotopt2_wave_file}" == "" ]]; then
  echo "--wave-file flag is required"
  exit 1
fi

if [[ ! -x "$(command -v ${gotopt2_viewer_binary})" ]]; then
  echo "could not find viewer binary: ${gotopt2_viwer_binary}"
  exit 1
fi
file_to_view="${gotopt2_wave_file}"
if [[ ! -f "${file_to_view}" ]]; then
  echo "can not find file to view: ${file_to_view}"
  exit 1
fi

"${gotopt2_viewer_binary}" \
  ${gotopt2_args__[@]} \
  "${file_to_view}"
