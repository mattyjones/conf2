#include <stdio.h>
#include "c2io/c2io_stream.h"

// Bridge function to call C function pointer in a given language to convert from language
// data streams to go data streams
int c2io_read_stream(c2io_read_stream_impl impl_func, void *sink, void *resourceLoader, void *buffPtr, char *resourceId) {
    return (*impl_func)(sink, resourceLoader, buffPtr, resourceId);
}
