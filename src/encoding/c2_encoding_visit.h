#ifndef C2_ENCODING_VISIT_H
#define C2_ENCODING_VISIT_H

// Each language (python, java, etc.) will implement at most one of these function
// pointers to convert data streams.  You will pass your function pointer when calling
// restconf_NewInputStream
//
typedef int (*c2_encoding_visit_impl)();


#endif