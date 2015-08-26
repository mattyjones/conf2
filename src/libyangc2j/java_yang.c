#include "java_yang.h"

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

bool checkDriverError(JNIEnv *env, GoInterface *err) {
  if ((*env)->ExceptionCheck(env)) {
    jthrowable exception = (*env)->ExceptionOccurred(env);
    (*env)->ExceptionClear(env);

    char *msg = NULL;
    jclass driver_class = (*env)->FindClass(env, "org/conf2/yang/driver/Driver");
    jmethodID print_err = (*env)->GetStaticMethodID(env, driver_class, "printException", "(Ljava/lang/Throwable;)Ljava/lang/String;");
    if (!(*env)->ExceptionCheck(env)) {
      jobject err_message = (*env)->CallStaticObjectMethod(env, driver_class, print_err, exception);
      if (!(*env)->ExceptionCheck(env)) {
        msg = (char *)(*env)->GetStringUTFChars(env, err_message, 0);
      }
    }
    if (msg == NULL) {
      msg = get_exception_message(env, exception);
    }

    *err = yangc2_new_driver_error(msg);
    (*env)->ExceptionClear(env);
    return true;
  }

  return false;
}

char *get_exception_message(JNIEnv *env, jthrowable err) {
  jclass throwable_class = (*env)->GetObjectClass(env, err);
  jmethodID get_message = (*env)->GetMethodID(env, throwable_class, "getMessage", "()Ljava/lang/String;");
  jstring msg = (*env)->CallObjectMethod(env, err, get_message);
  if (msg != NULL) {
    return (char *)(*env)->GetStringUTFChars(env, msg, 0);
  }
  jmethodID to_string = (*env)->GetMethodID(env, throwable_class, "toString", "()Ljava/lang/String;");
  jstring as_string = (*env)->CallObjectMethod(env, err, to_string);
  return (char *)(*env)->GetStringUTFChars(env, as_string, 0);
}

void initJvmReference(JNIEnv* env) {
  // Hook - get reference to current VM
  if (jvm == NULL) {
    (*env)->GetJavaVM(env, &jvm);
  }
}

jobject makeDriverHandle(JNIEnv *env, GoInterface iface) {
  GoInt err;
  jclass driverHandleCls = (*env)->FindClass(env, "org/conf2/yang/driver/DriverHandle");
  if (!(err = checkError(env))) {
    jmethodID driverHandleCtor = (*env)->GetMethodID(env, driverHandleCls, "<init>", "([B)V");
    if (!(err = checkError(env))) {
      jobject jbuffer = (*env)->NewByteArray(env, sizeof(iface));
      void* cbuffer = (*env)->GetByteArrayElements(env, jbuffer, 0);
      memcpy(cbuffer, &iface, sizeof(iface));
      (*env)->SetByteArrayRegion(env, jbuffer, 0, sizeof(iface), cbuffer);
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
  jclass driverHandleCls = (*env)->FindClass(env, "org/conf2/yang/driver/DriverHandle");
  if (!(err = checkError(env))) {
    jfieldID ifaceField = (*env)->GetFieldID(env, driverHandleCls, "reference", "Ljava/nio/ByteBuffer;");
    if (!(err = checkError(env))) {
      jbyteArray ref = (jbyteArray) (*env)->GetObjectField(env, driverHandle, ifaceField);
      if (!(err = checkError(env))) {
        void *ifaceBytes = (*env)->GetDirectBufferAddress(env, ref);
        if (ifaceBytes != NULL) {
          memcpy(iface, ifaceBytes, sizeof(*iface));
        } else {
          err = RC_BAD;
        }
      }
    }
  }
}

JNIEnv *getCurrentJniEnv() {
  JNIEnv* env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&env, NULL);
  return env;
}

JNIEXPORT void JNICALL Java_org_conf2_yang_driver_Driver_initializeDriver
  (JNIEnv *env, jobject jobj) {
  initJvmReference(env);
}

