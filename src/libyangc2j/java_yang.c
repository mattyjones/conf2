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

yangc2j_method get_static_driver_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  yangc2j_method m;
  m.methodId = NULL;
  m.cls = NULL;
  m.cls = (*env)->FindClass(env, "org/conf2/yang/driver/Driver");
  if (checkDriverError(env, err)) {
    return m;
  }
  m.methodId = (*env)->GetStaticMethodID(env, m.cls, method_name, signature);
  if (checkDriverError(env, err)) {
    return m;
  }
  return m;
}

bool checkDriverError(JNIEnv *env, GoInterface *err) {
  if ((*env)->ExceptionCheck(env)) {
    jthrowable exception = (*env)->ExceptionOccurred(env);
    (*env)->ExceptionClear(env);

    char *msg = NULL;
    yangc2j_method print_err = get_static_driver_method(env, err, "printException", "(Ljava/lang/Throwable;)Ljava/lang/String;");
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

void yangc2j_release_x_array(void *x_array_ref, void *errPtr) {
  yangc2j_array *x_array = (yangc2j_array *)x_array_ref;
  JNIEnv* env = getCurrentJniEnv();
  (*env)->DeleteGlobalRef(env, x_array->j_buff);
  if (x_array->cstr_list != NULL) {
    free(x_array->cstr_list);
  }
  free(x_array_ref);
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

yangc2j_array* yangc2j_new_cstr_list(JNIEnv *env, jobjectArray j_str_array, GoInterface *err) {
  yangc2j_array* list = (yangc2j_array*) calloc(1, sizeof(yangc2j_array));
  list->list_len = (*env)->GetArrayLength(env, j_str_array);
  if (checkDriverError(env, err)) {
    return NULL;
  }

  yangc2j_method to_buff = get_static_driver_method(env, err, "encodeCStrArray", "([Ljava/lang/String;)Ljava/nio/ByteBuffer;");
  if (to_buff.methodId == NULL) {
    free(list);
    return NULL;
  }
  jobject j_buff = (*env)->CallStaticObjectMethod(env, to_buff.cls, to_buff.methodId, j_str_array);
  if (checkDriverError(env, err)) {
    free(list);
    return NULL;
  }

  int bufflen = (*env)->GetDirectBufferCapacity(env, j_buff);
printf("java_yang.c: bufflen=%d\n", bufflen);
  list->j_buff = (*env)->NewGlobalRef(env, j_buff);
  list->handle = yangc2_handle_new((void *)list, &yangc2j_release_x_array);
  list->cstr_list = (char **)calloc(list->list_len, sizeof(char *));
  char *buff = (char *)(*env)->GetDirectBufferAddress(env, j_buff);
  int i;
  int j = 0;
  for (i = 0; j < bufflen && i < list->list_len; i++) {
    list->cstr_list[i] = &buff[j];
printf("java_yang.c: s[%d]=%s\n", i, list->cstr_list[i]);
    for (;buff[j] != 0 && j < bufflen; j++) {
    }
    j++; // null term
  }

  return list;
}

yangc2j_array* yangc2j_new_bool_list(JNIEnv *env, jobjectArray j_bool_array, GoInterface *err) {
  yangc2j_array* list = (yangc2j_array*) calloc(1, sizeof(yangc2j_array));
  list->list_len = (*env)->GetArrayLength(env, j_bool_array);
  if (checkDriverError(env, err)) {
    free(list);
    return NULL;
  }

  yangc2j_method to_buff = get_static_driver_method(env, err, "encodeCBoolArray", "([Z)Ljava/nio/ByteBuffer;");
  if (to_buff.methodId == NULL) {
    free(list);
    return NULL;
  }
  jobject j_buff = (*env)->CallStaticObjectMethod(env, to_buff.cls, to_buff.methodId, j_bool_array);
  if (checkDriverError(env, err)) {
    free(list);
    return NULL;
  }
  list->j_buff = (*env)->NewGlobalRef(env, j_buff);
  list->handle = yangc2_handle_new((void *)list, &yangc2j_release_x_array);
  list->bool_list = (short *)(*env)->GetDirectBufferAddress(env, j_buff);
  return list;
}

yangc2j_array* yangc2j_new_int_list(JNIEnv *env, jobjectArray j_int_array, GoInterface *err) {
  yangc2j_array* list = (yangc2j_array*) calloc(1, sizeof(yangc2j_array));
  list->list_len = (*env)->GetArrayLength(env, j_int_array);
  if (checkDriverError(env, err)) {
    free(list);
    return NULL;
  }

  yangc2j_method to_buff = get_static_driver_method(env, err, "encodeCIntArray", "([I)Ljava/nio/ByteBuffer;");
  if (to_buff.methodId == NULL) {
    free(list);
    return NULL;
  }
  jobject j_buff = (*env)->CallStaticObjectMethod(env, to_buff.cls, to_buff.methodId, j_int_array);
  if (checkDriverError(env, err)) {
    free(list);
    return NULL;
  }
  list->j_buff = (*env)->NewGlobalRef(env, j_buff);
  list->handle = yangc2_handle_new((void *)list, &yangc2j_release_x_array);
  list->int_list = (int *)(*env)->GetDirectBufferAddress(env, j_buff);
  return list;
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