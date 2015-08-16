#include <string.h>
#include <stdbool.h>
#include "libyangc2j/libyangc2j.h"
#include "org_conf2_yang_driver_Driver.h"
#include "yang/driver/yangc2_stream.h"
#include "driver.h"

bool checkFsError(JNIEnv *env, GoInterface *err) {
  if ((*env)->ExceptionCheck(env)) {
    jthrowable exception = (*env)->ExceptionOccurred(env);
    char *msg = get_exception_message(env, exception);
    *err = yangc2_new_fs_error(msg);
    (*env)->ExceptionClear(env);
    return true;
  }

  return false;
}

void *java_open_stream(void *source_handle, char *resId, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();
  jclass loaderIface = (*env)->FindClass(env, "org/conf2/yang/comm/DataSource");
  if (checkFsError(env, err)) {
    return NULL;
  }
  jmethodID getResourceMethod = (*env)->GetMethodID(env, loaderIface, "getResource", "(Ljava/lang/String;)Ljava/io/InputStream;");
  if (checkFsError(env, err)) {
    return NULL;
  }
  jclass inputStreamCls = (*env)->FindClass(env, "java/io/InputStream");
  if (checkFsError(env, err)) {
    return NULL;
  }
  jobject resourceIdStr = (*env)->NewStringUTF(env, resId);
  jobject inputStream = (*env)->CallObjectMethod(env, source_handle, getResourceMethod, resourceIdStr);
  if (checkFsError(env, err)) {
    return NULL;
  }
  return inputStream;
}

int java_read_stream(void *stream_handle, void *buffSlicePtr, int maxAmount, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();
  GoSlice buff = *((GoSlice *)buffSlicePtr);
  jobject inputStream = stream_handle;
  jclass inputStreamCls = (*env)->FindClass(env, "java/io/InputStream");
  if (checkFsError(env, err)) {
    return 0;
  }
  jmethodID readMethod = (*env)->GetMethodID(env, inputStreamCls, "read", "([BII)I");
  if (checkFsError(env, err)) {
    return 0;
  }
  // TODO: for performance, reuse byte buffer between reads. ideally figure out how to read
  // straight into given buffer w/o allocating but i couldn't figure that out
  jobject buffer = (*env)->NewByteArray(env, buff.cap);
  jint amountRead = (*env)->CallIntMethod(env, inputStream, readMethod, buffer, 0, buff.cap);
  if (checkFsError(env, err)) {
    return 0;
  }
  if (amountRead > 0) {
    void* chunk = (*env)->GetByteArrayElements(env, buffer, 0);
    if (chunk == NULL) {
      *err = yangc2_new_fs_error("Could not allocate java byte buffer");
    } else {
      memcpy(buff.data, chunk, amountRead);
      //buff.len = amountRead;
    }
  }
  (*env)->DeleteLocalRef(env, buffer);
  return amountRead;
}

void java_close_stream(void *stream_handle, void *errPtr) {
  GoInterface *err = (GoInterface *) errPtr;
  JNIEnv* env = getCurrentJniEnv();
  jobject inputStream = stream_handle;
  jclass inputStreamCls = (*env)->FindClass(env, "java/io/InputStream");
  if (checkFsError(env, err)) {
    return;
  }
  jmethodID closeMethod = (*env)->GetMethodID(env, inputStreamCls, "close", "()V");
  (*env)->CallObjectMethod(env, inputStream, closeMethod);
  checkFsError(env, err);
}

JNIEXPORT jstring JNICALL Java_org_conf2_yang_comm_Driver_echoTest
  (JNIEnv *env, jclass serviceCls, jobject resourceLoader, jstring resourceId) {
    initJvmReference(env);
    GoInterface source = yangc2_new_driver_resource_source(&java_open_stream, &java_read_stream, &java_close_stream, resourceLoader);
    const char *cResourceId = (*env)->GetStringUTFChars(env, resourceId, 0);
    char *results = yangc2_echo_test(source, (char *)cResourceId);
    (*env)->ReleaseStringUTFChars(env, resourceId, cResourceId);
    return (*env)->NewStringUTF(env, results);
}

JNIEXPORT jobject JNICALL Java_org_conf2_yang_comm_Driver_newDataSource
  (JNIEnv *env, jobject driver, jobject dataSource) {
  GoInterface ds = yangc2_new_driver_resource_source(&java_open_stream, &java_read_stream, &java_close_stream, dataSource);
  jobject dsHandle = makeDriverHandle(env, ds);
  return dsHandle;
}
