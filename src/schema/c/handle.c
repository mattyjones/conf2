#include <stdio.h>
#include "conf2/handle.h"

// Bridge functions to call C function pointer in a given language drivers

void conf2_handle_release_bridge(conf2_handle_release_impl impl_func, void *handle, void *errPtr) {
    (*impl_func)(handle, errPtr);
}
