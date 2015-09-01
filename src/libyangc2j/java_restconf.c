#include "yang-c2/java.h"
#include "org_conf2_restconf_Service.h"
#include "restconf-c2/driver.h"


JNIEXPORT void JNICALL Java_org_conf2_restconf_Service_sendSetDocRoot
  (JNIEnv *env, jobject j, jlong service_hnd_id, jobject j_stream_source) {

  jobject j_g_stream_source = (*env)->NewGlobalRef(env, j_stream_source);
  void *stream_source_hnd_id = yangc2_handle_new(j_g_stream_source, &java_release_global_ref);
  void *go_stream_source_hnd_id = yangc2_new_driver_resource_source(stream_source_hnd_id, &java_open_stream,
    &java_read_stream);

  restconfc2_set_doc_root((void *)service_hnd_id, go_stream_source_hnd_id);
}

JNIEXPORT jlong JNICALL Java_org_conf2_restconf_Service_newService
  (JNIEnv *env, jobject j) {
  return (jlong)restconfc2_service_new();
}

JNIEXPORT void JNICALL Java_org_conf2_restconf_Service_startService
  (JNIEnv *env, jobject j, jlong service_hnd_id) {
  restconfc2_service_start((void *)service_hnd_id);
}

JNIEXPORT void JNICALL Java_org_conf2_restconf_Service_registerBrowserWithService
  (JNIEnv *env, jobject service, jlong service_hnd_id, jlong module_hnd_id, jobject j_browser) {

  jobject j_g_browser = (*env)->NewGlobalRef(env, j_browser);
  void *browser_hnd_id = yangc2_handle_new(j_g_browser, &java_release_global_ref);

  void *go_browser_hnd_id = yangc2_new_browser(
        &java_browse_enter,
        &java_browse_iterate,
        &java_browse_read,
        &java_browse_edit,
        &java_browse_choose,
        &java_browse_exit,
        &java_browse_root_selector,
        (void *)module_hnd_id,
        browser_hnd_id);

  // TODO: Error handle
  restconfc2_register_browser((void *)service_hnd_id, go_browser_hnd_id);
}