#include <stdlib.h>
#include <stdio.h>
#include "yang-c2/java.h"

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

yangc2j_method yangc2j_static_method(JNIEnv *env, GoInterface *err, char *class_name, char *method_name, char *signature) {
  yangc2j_method m;
  m.methodId = NULL;
  m.cls = NULL;
  m.cls = (*env)->FindClass(env, class_name);
  if (checkDriverError(env, err)) {
    return m;
  }
  m.methodId = (*env)->GetStaticMethodID(env, m.cls, method_name, signature);
  if (checkDriverError(env, err)) {
    return m;
  }
  return m;
}

yangc2j_method yangc2j_static_driver_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  return yangc2j_static_method(env, err, "org/conf2/yang/driver/Driver", method_name, signature);
}

bool checkDriverError(JNIEnv *env, GoInterface *err) {
  if ((*env)->ExceptionCheck(env)) {
    jthrowable exception = (*env)->ExceptionOccurred(env);
    (*env)->ExceptionClear(env);

    char *msg = NULL;
    yangc2j_method print_err = yangc2j_static_driver_method(env, err, "printException", "(Ljava/lang/Throwable;)Ljava/lang/String;");
    if (print_err.methodId != NULL) {
      jobject err_message = (*env)->CallStaticObjectMethod(env, print_err.cls, print_err.methodId, exception);
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

void yangc2j_release_cstr(void *cstr_ref, void *errPtr) {
  yangc2j_cstr *chars = (yangc2j_cstr *)cstr_ref;
  JNIEnv* env = getCurrentJniEnv();
  (*env)->ReleaseStringUTFChars(env, chars->j_string, chars->cstr);
  free(cstr_ref);
}

yangc2j_cstr* yangc2j_new_cstr(jobject j_string) {
  JNIEnv* env = getCurrentJniEnv();
  yangc2j_cstr* ref = (yangc2j_cstr* )malloc(sizeof(yangc2j_cstr));
  ref->j_string = j_string;
  ref->cstr = (*env)->GetStringUTFChars(env, j_string, 0);
  ref->handle = yangc2_handle_new(ref, yangc2j_release_cstr);
  return ref;
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

void yangc2j_release_global_ref(void *ref, void *errPtr) {
  JNIEnv* env = getCurrentJniEnv();
  jobject j_ref = (jobject)ref;
  (*env)->DeleteGlobalRef(env, j_ref);
}

JNIEXPORT void JNICALL Java_org_conf2_yang_driver_Driver_releaseHandle
  (JNIEnv *env, jobject driver, jlong hnd_id) {
  yangc2_handle_release((void *)hnd_id);
}