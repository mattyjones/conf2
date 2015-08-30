#ifndef YANGC2_JAVA_H
#define YANGC2_JAVA_H

#include <string.h>
#include <stdbool.h>
#include <jni.h>
#include "yang-c2/driver.h"
#include "yang-c2/browse.h"

typedef enum {
  RC_OK = 0,
  RC_BAD = -1
} RcError;

typedef struct _yangc2j_method {
  jobject cls;
  jmethodID methodId;
} yangc2j_method;

RcError checkError(JNIEnv *env);

bool checkDriverError(JNIEnv *env, GoInterface *err);

char *get_exception_message(JNIEnv *env, jthrowable err);

void initJvmReference(JNIEnv* env);

// Wrap java object in releasable go handle by putting java object in global
// heap and creating a ApiHandle in Go that should be released when Go is done
// with object
//
// Example:
//   jobject j_g_obj = (*env)->NewGlobalRef(env, j_obj);
//   void *api_handle = yangc2_new_api_handle(j_g_obj, &java_release_global_ref);
//
void java_release_global_ref(void *ref, void *errPtr);

typedef struct _java_string_chars {
    void *handle;
    jobject j_string;
    const char* chars;
} java_string_chars;

void java_release_string_chars(void *chars, void *errPtr);
java_string_chars* java_new_string_chars(jobject j_string);
int java_read_stream(void *stream_handle, void *buffSlicePtr, int maxAmount, void *errPtr);
void *java_open_stream(void *source_handle, char *resId, void *errPtr);
void java_close_stream(void *stream_handle, void *errPtr);

JNIEnv *getCurrentJniEnv();

#endif