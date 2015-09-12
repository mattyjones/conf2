#include <stdio.h>
#include "conf2/stream.h"

// Bridge functions to call C function pointer in a given language to convert from language
// data streams to go data streams
void *conf2_open_stream(conf2_open_stream_impl impl_func, void *source_handle, char *resId, void *fsErr) {
    return (*impl_func)(source_handle, resId, fsErr);
}

int conf2_read_stream(conf2_read_stream_impl impl_func, void *stream_handle, void *buffPtr, int maxAmount, void *fs_err) {
    return (*impl_func)(stream_handle, buffPtr, maxAmount, fs_err);
}
