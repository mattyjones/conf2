#ifndef YANGC2J_H
#define YANGC2J_H

#include <string.h>
#include <stdbool.h>
#include <jni.h>
#include "yang/driver.h"
#include "yang/driver/yangc2_browse.h"

typedef enum {
  RC_OK = 0,
  RC_BAD = -1
} RcError;

typedef struct _yangc2j_method {
  jobject cls;
  jmethodID methodId;
} yangc2j_method;

RcError checkError(JNIEnv *env);
bool checkDriverError(JNIEnv *env, GoInterface *err);

char *get_exception_message(JNIEnv *env, jthrowable err);

void initJvmReference(JNIEnv* env);

jobject makeDriverHandle(JNIEnv *env, GoInterface iface);

RcError resolveDriverHandle(JNIEnv *env, jobject driverHandle, GoInterface *iface);

JNIEnv *getCurrentJniEnv();

#endif