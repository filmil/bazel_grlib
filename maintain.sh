#!/bin/bash
set -e

OTHER_ARGS=()
while [[ $# -gt 0 ]]; do
  case $1 in
    --maintain_bin=*)
      MAINTAIN_BIN="${1#*=}"
      shift
      ;;
    --gen_build_files=*)
      GEN_BUILD_FILES="${1#*=}"
      shift
      ;;
    --in2kconfig=*)
      IN2KCONFIG="${1#*=}"
      shift
      ;;
    --python_bin=*)
      PYTHON_BIN="${1#*=}"
      shift
      ;;
    *)
      OTHER_ARGS+=("$1")
      shift
      ;;
  esac
done

# Find GRLIB_SRCS. In runfiles it is typically at external/+http_archive+grlib_srcs
GRLIB_SRCS=""
if [ -d "external/+http_archive+grlib_srcs" ]; then
    GRLIB_SRCS=$(readlink -f "external/+http_archive+grlib_srcs")
elif [ -d "../+http_archive+grlib_srcs" ]; then
    GRLIB_SRCS=$(readlink -f "../+http_archive+grlib_srcs")
fi

if [ -z "$GRLIB_SRCS" ]; then
    # Try finding it recursively in the current directory (runfiles root)
    GRLIB_SRCS=$(find . -name "lib" -type d -path "*/external/*" | head -n 1 | xargs dirname | xargs readlink -f)
fi

export GRLIB_SRCS
if [ -z "$GRLIB_SRCS" ]; then
    echo "ERROR: Could not find GRLIB_SRCS"
    exit 1
fi

# Ensure paths are absolute if they are relative (relative to runfiles root)
[[ "$MAINTAIN_BIN" != /* ]] && MAINTAIN_BIN="$(pwd)/$MAINTAIN_BIN"
[[ "$GEN_BUILD_FILES" != /* ]] && GEN_BUILD_FILES="$(pwd)/$GEN_BUILD_FILES"
[[ "$IN2KCONFIG" != /* ]] && IN2KCONFIG="$(pwd)/$IN2KCONFIG"

exec "$MAINTAIN_BIN" \
    --gen_build_files "$GEN_BUILD_FILES" \
    --in2kconfig "$IN2KCONFIG" \
    "${OTHER_ARGS[@]}"
