#ifndef YANGC2_CONTEXT_H
#define YANGC2_CONTEXT_H

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to release resources.  You will pass your function pointer when calling
// yangc2_handle_new
//
typedef void (*yangc2_handle_release_impl)(void *handle, void *errPtr);

#endif
