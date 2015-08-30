#ifndef YANGC2_STREAM_H
#define YANGC2_STREAM_H

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to convert data streams.  You will pass your function pointer when calling
// yangc2_new_driver_resource_source
//
typedef void *(*yangc2_open_stream_impl)(void *source_handle, char *resId, void *fsErr);
typedef int (*yangc2_read_stream_impl)(void *stream_handle, void *buffPtr, int maxAmout, void *fs_err);

#endif
