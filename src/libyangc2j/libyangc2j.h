#ifndef YANGC2J_H
#define YANGC2J_H

#include <string.h>
#include <stdbool.h>
#include <jni.h>

typedef enum {
  RC_OK = 0,
  RC_BAD = -1
} RcError;

RcError checkError(JNIEnv *env);

char *get_exception_message(JNIEnv *env, jthrowable err);

void initJvmReference(JNIEnv* env);

jobject makeDriverHandle(JNIEnv *env, GoInterface iface);

RcError resolveDriverHandle(JNIEnv *env, jobject driverHandle, GoInterface *iface);

#endif