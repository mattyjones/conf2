/* DO NOT EDIT THIS FILE - it is machine generated */
#include <jni.h>
/* Header for class org_conf2_yang_comm_Driver */

#ifndef _Included_org_conf2_yang_comm_Driver
#define _Included_org_conf2_yang_comm_Driver
#ifdef __cplusplus
extern "C" {
#endif
/*
 * Class:     org_conf2_yang_comm_Driver
 * Method:    initializeDriver
 * Signature: ()V
 */
JNIEXPORT void JNICALL Java_org_conf2_yang_comm_Driver_initializeDriver
  (JNIEnv *, jobject);

/*
 * Class:     org_conf2_yang_comm_Driver
 * Method:    echoTest
 * Signature: (Lorg/conf2/yang/comm/DataSource;Ljava/lang/String;)Ljava/lang/String;
 */
JNIEXPORT jstring JNICALL Java_org_conf2_yang_comm_Driver_echoTest
  (JNIEnv *, jclass, jobject, jstring);

/*
 * Class:     org_conf2_yang_comm_Driver
 * Method:    newDataSource
 * Signature: (Lorg/conf2/yang/comm/DataSource;)Lorg/conf2/yang/comm/DriverHandle;
 */
JNIEXPORT jobject JNICALL Java_org_conf2_yang_comm_Driver_newDataSource
  (JNIEnv *, jobject, jobject);

#ifdef __cplusplus
}
#endif
#endif
