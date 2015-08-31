#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include "yang-c2/java.h"
#include "org_conf2_yang_driver_Driver.h"
#include "yang-c2/stream.h"
#include "yang-c2/driver.h"

void java_close_stream(void *stream_handle, void *errPtr) {
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
  java_release_global_ref(stream_handle, errPtr);
}

void *java_open_stream(void *source_handle, char *resId, void *errPtr) {
printf("java_stream.c:java_open_stream source_handle=%p\n", source_handle);
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
    *err = yangc2_new_driver_error("Stream not found");
  }
  jobject j_g_inputstream = (*env)->NewGlobalRef(env, j_inputstream);
  void *handle = yangc2_handle_new(j_g_inputstream, &java_close_stream);
printf("java_stream.c:java_open_stream LEAVIG handle=%p, j_g_inputstream=%p\n", handle, j_g_inputstream);
  return handle;
}

int java_read_stream(void *stream_handle, void *buffSlicePtr, int maxAmount, void *errPtr) {
printf("java_stream.c:java_read_stream stream_handle=%p\n", stream_handle);
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
      *err = yangc2_new_driver_error("Could not allocate java byte buffer");
    } else {
      memcpy(buff.data, chunk, amountRead);
      //buff.len = amountRead;
    }
  }
  (*env)->DeleteLocalRef(env, buffer);
  return amountRead;
}

