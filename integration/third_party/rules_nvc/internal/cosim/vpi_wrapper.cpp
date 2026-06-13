#include "vhpi_user.h"
#ifdef VPI_BINDINGS_HEADER
#include VPI_BINDINGS_HEADER
#else
#include "vpi_bindings.hpp"
#endif
#include <unordered_map>
#include <iostream>
#include <string>
#include <cstdio>
#include <stdint.h>
#include <dlfcn.h>

std::unordered_map<int, InstanceState> instances;

extern "C" {

#ifndef VERILATOR_STEP_CALL
#define VERILATOR_STEP_CALL verilator_step_call
#endif

void VERILATOR_STEP_CALL(int32_t id) {
    if (instances.find(id) == instances.end()) {
        std::string path_prefix = "";

#ifdef PATH_PREFIX
        path_prefix = PATH_PREFIX;
#else
        path_prefix = ":top_tb:dut_inst"; // Fallback
#endif

        init_bindings(id, instances[id], path_prefix);
    }

    InstanceState& state = instances[id];

    sync_inputs(state);
    eval_model(state);
    sync_outputs(state);
}

void (*vhpi_startup_routines[])(void) = {
    nullptr
};

} // extern "C"

