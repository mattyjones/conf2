#ifndef YANGC2_COMM_STREAM_H
#define YANGC2_COMM_STREAM_H

// Each language (python, java, etc.) will implement at most one of these function
// pointers to convert data streams.  You will pass your function pointer when calling
// restconf_NewInputStream
//
typedef int (*yangc2_comm_read_stream_impl)(void *sink, void *rloader, void *buffPtr, char *resId);

#endif
