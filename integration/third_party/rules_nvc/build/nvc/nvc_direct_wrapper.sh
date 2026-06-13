#!/bin/bash
set -e
_ld_so=""
IFS=':' read -ra LD_DIRS <<< "${NVC_LD_LIBRARY_PATH:-}"
for _dir in "${LD_DIRS[@]}"; do
  if [[ -x "$_dir/ld-linux-x86-64.so.2" ]]; then
    _ld_so="$_dir/ld-linux-x86-64.so.2"
    break
  fi
done
export LD_LIBRARY_PATH="${NVC_LD_LIBRARY_PATH:-$LD_LIBRARY_PATH}"
exec ${_ld_so:+"$_ld_so"} "$@"
