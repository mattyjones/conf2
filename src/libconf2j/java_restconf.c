#include "conf2/java.h"
#include "org_conf2_restconf_Service.h"
#include "conf2/restconf.h"


JNIEXPORT void JNICALL Java_org_conf2_restconf_Service_sendSetDocRoot
  (JNIEnv *env, jobject j, jlong service_hnd_id, jobject j_stream_source) {

    jobject j_g_stream_source = (*env)->NewGlobalRef(env, j_stream_source);
    void *stream_source_hnd_id = conf2_handle_new(j_g_stream_source, &conf2j_release_global_ref);
    void *go_stream_source_hnd_id = conf2_new_driver_resource_source(stream_source_hnd_id, &conf2j_open_stream,
        &conf2j_read_stream);

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

    void *go_browser_hnd_id = conf2j_browse_new(env, module_hnd_id, j_browser);

    // TODO: Error handle
    restconfc2_register_browser((void *)service_hnd_id, go_browser_hnd_id);
}
