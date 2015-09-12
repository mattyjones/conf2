#ifndef CONF2_STREAM_H
#define CONF2_STREAM_H

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to convert data streams.  You will pass your function pointer when calling
// conf2_new_driver_resource_source
//
typedef void *(*conf2_open_stream_impl)(void *source_handle, char *resId, void *fsErr);
typedef int (*conf2_read_stream_impl)(void *stream_handle, void *buffPtr, int maxAmount, void *fs_err);

#endif
