#include "libyangc2j.h"

JavaVM* jvm = NULL;

/**
 * PERFORMANCE TODO: Cache all java reflection calls
 */

RcError checkError(JNIEnv *env) {
  if ((*env)->ExceptionCheck(env)) {
    (*env)->ExceptionDescribe(env);
    (*env)->ExceptionClear(env);
    return RC_BAD;
  }
  return RC_OK;
}

char *get_exception_message(JNIEnv *env, jthrowable err) {
  jclass throwable_class = (*env)->FindClass(env, "java/lang/Throwable");
  jmethodID to_string = (*env)->GetMethodID(env, throwable_class, "toString", "()Ljava/lang/String;");
  jstring msg = (*env)->CallObjectMethod(env, throwable_class, to_string);
  return (char *)(*env)->GetStringUTFChars(env, msg, 0);
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

JNIEXPORT void JNICALL Java_org_conf2_yang_comm_Driver_initializeDriver
  (JNIEnv *env, jobject jobj) {
  initJvmReference(env);
}


JNIEXPORT jobject JNICALL Java_org_conf2_yang_Loader_loadModule
  (JNIEnv *env, jclass loaderClass, jobject dataSourceHandle, jstring resource) {
  GoInterface dsIface;
  resolveDriverHandle(env, dataSourceHandle, &dsIface);
  char *resourceStr = (char *)(*env)->GetStringUTFChars(env, resource, 0);
  GoInterface err = yangc2_load_module_from_resource_source(dsIface, resourceStr);
  return NULL;
}

