#ifndef CONF2_JAVA_H
#define CONF2_JAVA_H

#include <string.h>
#include <stdbool.h>
#include <jni.h>
#include "conf2/schema.h"
#include "conf2/browse.h"

typedef enum {
  RC_OK = 0,
  RC_BAD = -1
} RcError;

typedef struct _conf2j_method {
  jobject cls;
  jmethodID methodId;
} conf2j_method;

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
//   void *api_handle = conf2_new_api_handle(j_g_obj, &conf2j_release_global_ref);
//
void conf2j_release_global_ref(void *ref, void *errPtr);

typedef struct _conf2j_cstr {
    void *handle;
    jobject j_string;
    const char* cstr;
} conf2j_cstr;

void conf2j_release_cstr(void *chars, void *errPtr);
conf2j_cstr* conf2j_new_cstr(jobject j_string);

int conf2j_read_stream(void *stream_handle, void *buffSlicePtr, int maxAmount, void *errPtr);
void *conf2j_open_stream(void *source_handle, char *resId, void *errPtr);
void conf2j_close_stream(void *stream_handle, void *errPtr);

void* conf2j_browse_root_selector(void *browser_handle, void *browse_err);
void* conf2j_browse_enter(void *selection_handle, char *ident, short *found, void *browse_err);
short conf2j_browse_iterate(void *selection_handle, void *key_data, int key_data_len, short first, void *browse_err);
void *conf2j_browse_read(void *selection_handle, char *ident, void **val_data_ptr, int *val_data_len_ptr, void *browse_err);
void conf2j_browse_edit(void *selection_handle, char *ident, int op, void *val_data, int val_data_len, void *browse_err);
char *conf2j_browse_choose(void *selection_handle, char *ident, void *browse_err);
void conf2j_browse_exit(void *selection_handle, char *ident, void *browse_err);
void *conf2j_browse_new(JNIEnv *env, jlong module_hnd_id, jobject j_browser);

conf2j_method conf2j_static_method(JNIEnv *env, GoInterface *err, char *class_name, char *method_name, char *signature);
conf2j_method conf2j_static_driver_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature);

#endif