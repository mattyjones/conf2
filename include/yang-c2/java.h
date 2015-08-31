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

JNIEnv *getCurrentJniEnv();

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

void* java_browse_root_selector(void *browser_handle, void *browse_err);
void* java_browse_enter(void *selection_handle, char *ident, short *found, void *browse_err);
short java_browse_iterate(void *selection_handle, char *encodedKeys, short first, void *browse_err);
void java_browse_read(void *selection_handle, char *ident, struct yangc2_browse_value *val, void *browse_err);
void java_browse_edit(void *selection_handle, char *ident, int op, struct yangc2_browse_value *val, void *browse_err);
char *java_browse_choose(void *selection_handle, char *ident, void *browse_err);
void java_browse_exit(void *selection_handle, char *ident, void *browse_err);


#endif