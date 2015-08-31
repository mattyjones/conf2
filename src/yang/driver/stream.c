#include <stdio.h>
#include "yang-c2/stream.h"

// Bridge functions to call C function pointer in a given language to convert from language
// data streams to go data streams
void *yangc2_open_stream(yangc2_open_stream_impl impl_func, void *source_handle, char *resId, void *fsErr) {
    return (*impl_func)(source_handle, resId, fsErr);
}

int yangc2_read_stream(yangc2_read_stream_impl impl_func, void *stream_handle, void *buffPtr, int maxAmount, void *fs_err) {
    return (*impl_func)(stream_handle, buffPtr, maxAmount, fs_err);
}
