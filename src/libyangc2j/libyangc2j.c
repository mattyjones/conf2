#include <string.h>
#include "org_conf2_yang_comm_Driver.h"
#include "yang/comm/yangc2_comm_stream.h"
#include "yang/comm.h"

JavaVM* jvm = NULL;

/**
 * PERFORMANCE: Cache all java reflaction calls
 */
typedef enum {
  RC_OK = 0,
  RC_BAD = -1
} RcError;

RcError checkError(JNIEnv *env) {
  if ((*env)->ExceptionCheck(env)) {
    (*env)->ExceptionDescribe(env);
    (*env)->ExceptionClear(env);
    return RC_BAD;
  }
  return RC_OK;
}

void initJvmReference(JNIEnv* env) {
  // Hook - get reference to current VM
  if (jvm == NULL) {
    (*env)->GetJavaVM(env, &jvm);
  }
}

jobject makeDriverHandle(JNIEnv *env, GoInterface iface) {
  GoInt err;
  jclass driverHandleCls = (*env)->FindClass(env, "org/conf2/yang/comm/DriverHandle");
  if (!(err = checkError(env))) {
    jmethodID driverHandleCtor = (*env)->GetMethodID(env, driverHandleCls, "<init>", "([B)V");
    if (!(err = checkError(env))) {
      jobject jbuffer = (*env)->NewByteArray(env, sizeof(iface));
      void* cbuffer = (*env)->GetByteArrayElements(env, jbuffer, 0);
      memcpy(cbuffer, &iface, sizeof(iface));
      jobject handle = (*env)->NewObject(env, driverHandleCls, driverHandleCtor, jbuffer);
      (*env)->ReleaseByteArrayElements(env, jbuffer, cbuffer, JNI_ABORT);
      (*env)->DeleteLocalRef(env, jbuffer);
      return handle;
    }
  }
  return NULL;
}

RcError resolveDriverHandle(JNIEnv *env, jobject driverHandle, GoInterface *iface) {
  GoInt err;
  jclass driverHandleCls = (*env)->FindClass(env, "org/conf2/yang/comm/DriverHandle");
  if (!(err = checkError(env))) {
    jfieldID ifaceField = (*env)->GetFieldID(env, driverHandleCls, "reference", "[B");
    if (!(err = checkError(env))) {
      jbyteArray ref = (jbyteArray) (*env)->GetObjectField(env, driverHandle, ifaceField);
      if (!(err = checkError(env))) {
        void *ifaceBytes = (*env)->GetDirectBufferAddress(env, ref);
        memcpy(iface, ifaceBytes, sizeof(*iface));
      }
    }
  }
}

JNIEnv *getCurrentJniEnv() {
  JNIEnv* env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&env, NULL);
  return env;
}

int java_read_stream(void *sinkPtr, void *resourceLoader, void *buffPtr, char *resourceId) {
  RcError err = RC_OK;
  GoInterface sink = *((GoInterface *) sinkPtr);
  GoSlice buff = *((GoSlice *)buffPtr);
  JNIEnv* env = getCurrentJniEnv();
  jclass loaderIface = (*env)->FindClass(env, "org/conf2/yang/comm/DataSource");
  if (err = checkError(env)) {
    return err;
  }
  jmethodID getResourceMethod = (*env)->GetMethodID(env, loaderIface, "getResource", "(Ljava/lang/String;)Ljava/io/InputStream;");
  if (err = checkError(env)) {
    return err;
  }
  jclass inputStreamCls = (*env)->FindClass(env, "java/io/InputStream");
  if (err = checkError(env)) {
    return err;
  }
  jobject resourceIdStr = (*env)->NewStringUTF(env, resourceId);
  jobject inputStream = (*env)->CallObjectMethod(env, resourceLoader, getResourceMethod, resourceIdStr);
  if (!(err = checkError(env))) {
    jmethodID readMethod = (*env)->GetMethodID(env, inputStreamCls, "read", "([BII)I");
    if (!(err = checkError(env))) {
      jobject buffer = (*env)->NewByteArray(env, buff.len);
      jint len = (*env)->CallIntMethod(env, inputStream, readMethod, buffer, 0, buff.len);
      void* chunk;
      if (!(err = checkError(env))) {
        while (len > 0) {
          chunk = (*env)->GetByteArrayElements(env, buffer, 0);
          if (chunk != NULL) {
            memcpy(buff.data, chunk, len);
            if (err = yangc2_comm_datasink_writedata(sink, buff)) {
              break;
            }
            (*env)->ReleaseByteArrayElements(env, buffer, chunk, JNI_ABORT);
            len = (*env)->CallIntMethod(env, inputStream, readMethod, buffer, 0, buff.len);
            if (err = checkError(env)) {
              break;
            }
          }
        }
      }
      (*env)->DeleteLocalRef(env, buffer);
    }
    jmethodID closeMethod = (*env)->GetMethodID(env, inputStreamCls, "close", "()V");
    (*env)->CallObjectMethod(env, inputStream, closeMethod);
    // closeQuietly - what would you do w/this error?
  }
  return err;
}

JNIEXPORT void JNICALL Java_org_conf2_yang_comm_Driver_initializeDriver
  (JNIEnv *env, jobject jobj) {
  initJvmReference(env);
}

JNIEXPORT jstring JNICALL Java_org_conf2_yang_comm_Driver_echoTest
  (JNIEnv *env, jclass serviceCls, jobject resourceLoader, jstring resourceId) {
    initJvmReference(env);
    GoInterface source = yangc2_comm_new_driver_data_source(&java_read_stream, resourceLoader);
    const char *cResourceId = (*env)->GetStringUTFChars(env, resourceId, 0);
    char *results = yangc2_comm_echo_test(source, (char *)cResourceId);
    (*env)->ReleaseStringUTFChars(env, resourceId, cResourceId);
    return (*env)->NewStringUTF(env, results);
}
