#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include "conf2/java.h"
#include "org_conf2_schema_driver_Driver.h"
#include "conf2/stream.h"
#include "conf2/schema.h"

void conf2j_close_stream(void *stream_handle, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();
  jobject inputStream = stream_handle;
  jclass inputStreamCls = (*env)->FindClass(env, "java/io/InputStream");
  if (checkDriverError(env, err)) {
    return;
  }
  jmethodID closeMethod = (*env)->GetMethodID(env, inputStreamCls, "close", "()V");
  (*env)->CallObjectMethod(env, inputStream, closeMethod);
  checkDriverError(env, err);
  conf2j_release_global_ref(stream_handle, errPtr);
}

void *conf2j_open_stream(void *source_handle, char *resId, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();

  jobject dataSource = (jobject)source_handle;
  jclass loaderIface = (*env)->GetObjectClass(env, dataSource);
  if (checkDriverError(env, err)) {
    return NULL;
  }
  jmethodID getResourceMethod = (*env)->GetMethodID(env, loaderIface, "getStream", "(Ljava/lang/String;)Ljava/io/InputStream;");
  if (checkDriverError(env, err)) {
    return NULL;
  }

  jobject resourceIdStr = (*env)->NewStringUTF(env, resId);
  jobject j_inputstream = (*env)->CallObjectMethod(env, source_handle, getResourceMethod, resourceIdStr);
  if (checkDriverError(env, err)) {
    return NULL;
  }
  if (j_inputstream == NULL) {
    *err = conf2_new_driver_error("Stream not found");
  }
  jobject j_g_inputstream = (*env)->NewGlobalRef(env, j_inputstream);
  void *handle = conf2_handle_new(j_g_inputstream, &conf2j_close_stream);
  return handle;
}

int conf2j_read_stream(void *stream_handle, void *buffSlicePtr, int maxAmount, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();
  GoSlice buff = *((GoSlice *)buffSlicePtr);
  jobject inputStream = stream_handle;
  jclass inputStreamCls = (*env)->GetObjectClass(env, inputStream);
  if (checkDriverError(env, err)) {
    return 0;
  }
  jmethodID readMethod = (*env)->GetMethodID(env, inputStreamCls, "read", "([BII)I");
  if (checkDriverError(env, err)) {
    return 0;
  }
  // TODO: for performance, reuse byte buffer between reads. ideally figure out how to read
  // straight into given buffer w/o allocating but i couldn't figure that out
  jobject buffer = (*env)->NewByteArray(env, buff.cap);
  jint amountRead = (*env)->CallIntMethod(env, inputStream, readMethod, buffer, 0, buff.cap);
  if (checkDriverError(env, err)) {
    return 0;
  }
  if (amountRead > 0) {
    void* chunk = (*env)->GetByteArrayElements(env, buffer, 0);
    if (chunk == NULL) {
      *err = conf2_new_driver_error("Could not allocate java byte buffer");
    } else {
      memcpy(buff.data, chunk, amountRead);
      //buff.len = amountRead;
    }
  }
  (*env)->DeleteLocalRef(env, buffer);
  return amountRead;
}

