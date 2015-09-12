#ifndef CONF2_CONTEXT_H
#define CONF2_CONTEXT_H

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to release resources.  You will pass your function pointer when calling
// conf2_handle_new
//
typedef void (*conf2_handle_release_impl)(void *handle, void *errPtr);

#endif
