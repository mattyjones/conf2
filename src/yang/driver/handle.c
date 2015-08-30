#include <stdio.h>
#include "yang-c2/handle.h"

// Bridge functions to call C function pointer in a given language drivers

void yangc2_handle_release_bridge(yangc2_handle_release_impl impl_func, void *handle, void *errPtr) {
    (*impl_func)(handle, errPtr);
}
