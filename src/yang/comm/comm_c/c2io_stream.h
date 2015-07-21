#ifndef C2IO_STREAM_H
#define C2IO_STREAM_H

// Each language (python, java, etc.) will implement at most one of these function
// pointers to convert data streams.  You will pass your function pointer when calling
// restconf_NewInputStream
//
typedef int (*c2io_read_stream_impl)(void *sink, void *rloader, void *buffPtr, char *resId);

#endif
