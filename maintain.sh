#!/bin/bash
set -e

# --- begin runfiles.bash initialization v3 ---
set -uo pipefail; f=bazel_tools/tools/bash/runfiles/runfiles.bash
if [ -z "${RUNFILES_DIR:-}" ] && [ -z "${RUNFILES_MANIFEST_FILE:-}" ]; then
  if [ -d "$0.runfiles" ]; then
    export RUNFILES_DIR="$0.runfiles"
  elif [ -f "$0.runfiles_manifest" ]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles_manifest"
  elif [ -f "$0.runfiles/MANIFEST" ]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles/MANIFEST"
  fi
fi
if [ -f "${RUNFILES_DIR:-/dev/null}/$f" ]; then
  source "${RUNFILES_DIR}/$f"
elif [ -f "${RUNFILES_MANIFEST_FILE:-/dev/null}" ]; then
  source "$(grep -sm1 "^$f " "$RUNFILES_MANIFEST_FILE" | cut -f2- -d' ')"
else
  rlocation() {
    if [ -f "$1" ]; then echo "$1";
    elif [ -f "external/$1" ]; then echo "external/$1";
    fi
  }
fi
# --- end runfiles.bash initialization v3 ---

# Locate components
GRLIB_SRCS_FILE=$(rlocation "grlib_srcs/lib/libs.txt")
if [ -n "$GRLIB_SRCS_FILE" ]; then
    GRLIB_SRCS=$(dirname "$(dirname "$GRLIB_SRCS_FILE")")
fi
[ -z "${GRLIB_SRCS:-}" ] && GRLIB_SRCS=$(find . -maxdepth 5 -name "lib" -type d | grep "grlib_srcs" | head -n 1 | xargs -r dirname)

export GRLIB_SRCS
if [ -z "$GRLIB_SRCS" ]; then
    echo "ERROR: Could not find GRLIB_SRCS."
    exit 1
fi

find_bin() {
    local name=$1
    local res=$(rlocation "grlib/third_party/grlib/scripts/$name")
    if [ -z "$res" ]; then
         res=$(rlocation "grlib/third_party/grlib/scripts/${name}_/${name}")
    fi
    if [ -z "$res" ]; then
        # Search runfiles tree
        res=$(find "${RUNFILES_DIR:-.}" -type f -executable -name "$name" | grep "scripts" | head -n 1)
    fi
    if [ -n "$res" ]; then
        readlink -f "$res"
    fi
}

BIN_maintain=$(find_bin maintain)
BIN_gen_build_files=$(find_bin gen_build_files)
BIN_in2kconfig=$(find_bin in2kconfig)
BIN_gen_master_kconfig=$(find_bin gen_master_kconfig)
BIN_gen_example_bazelrc=$(find_bin gen_example_bazelrc)
BIN_vhd_preprocess=$(find_bin vhd_preprocess)

if [ -z "$BIN_maintain" ]; then echo "ERROR: maintain not found"; exit 1; fi
if [ -z "$BIN_gen_build_files" ]; then echo "ERROR: gen_build_files not found"; exit 1; fi
if [ -z "$BIN_in2kconfig" ]; then echo "ERROR: in2kconfig not found"; exit 1; fi
if [ -z "$BIN_gen_master_kconfig" ]; then echo "ERROR: gen_master_kconfig not found"; exit 1; fi
if [ -z "$BIN_gen_example_bazelrc" ]; then echo "ERROR: gen_example_bazelrc not found"; exit 1; fi
if [ -z "$BIN_vhd_preprocess" ]; then echo "ERROR: vhd_preprocess not found"; exit 1; fi

# Change to workspace root in runfiles
WS_ROOT=$(dirname $(readlink -f $(rlocation grlib/Kconfig)))
cd "$WS_ROOT"

"$BIN_maintain" \
    --gen_build_files "$BIN_gen_build_files" \
    --in2kconfig "$BIN_in2kconfig" \
    --gen_master_kconfig "$BIN_gen_master_kconfig" \
    --vhd_preprocess "$BIN_vhd_preprocess" \
    "$@"

# Update example bazelrc if needed (only when running via 'bazel run')
if [[ ! " $@ " =~ " --check " ]]; then
    if [ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]; then
        echo "Updating docs/example.bazelrc..."
        cd "$BUILD_WORKSPACE_DIRECTORY"
        npx bazelisk query "@grlib_config//:all" --output=label_kind | grep "CONFIG_" > settings_list.txt
        "$BIN_gen_example_bazelrc" settings_list.txt docs/example.bazelrc
        rm settings_list.txt
    fi
fi
