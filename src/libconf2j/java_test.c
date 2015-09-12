#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include "conf2/java.h"
#include "org_conf2_schema_driver_DriverTestHarness.h"
#include "conf2/stream.h"
#include "conf2/schema.h"

JNIEXPORT jlong JNICALL Java_org_conf2_schema_driver_DriverTestHarness_newTestHarness
  (JNIEnv *env, jobject harness, jlong module_hnd_id, jobject j_browser) {
    void *go_browser_hnd_id = conf2j_browse_new(env, module_hnd_id, j_browser);
    return (jlong)conf2_testharness_new(go_browser_hnd_id);
}

JNIEXPORT jboolean JNICALL Java_org_conf2_schema_driver_DriverTestHarness_runTest
  (JNIEnv *env, jobject harness, jlong harness_hnd_id, jstring j_testname) {
    const char *c_testname = (*env)->GetStringUTFChars(env, j_testname, 0);
    jboolean passed = (jboolean) conf2_testharness_test_run((void *)harness_hnd_id, (char *)c_testname);
    (*env)->ReleaseStringUTFChars(env, j_testname, c_testname);
    return passed;
}

JNIEXPORT jstring JNICALL Java_org_conf2_schema_driver_DriverTestHarness_report
  (JNIEnv *env, jobject harness, jlong harness_hnd_id) {
    char *c_report = conf2_testharness_report((void *)harness_hnd_id);
    jobject j_report = (*env)->NewStringUTF(env, c_report);
    return j_report;
}
