#include <stdio.h>
#include "yangc2_comm_stream.h"

// Bridge function to call C function pointer in a given language to convert from language
// data streams to go data streams
int yangc2_comm_read_stream(yangc2_comm_read_stream_impl impl_func, void *sink, void *resourceLoader, void *buffPtr, char *resourceId) {
    return (*impl_func)(sink, resourceLoader, buffPtr, resourceId);
}
