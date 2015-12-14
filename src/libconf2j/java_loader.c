#include <string.h>
#include <stdbool.h>
#include <stdlib.h>
#include "conf2/java.h"
#include "org_conf2_schema_yang_ModuleLoader.h"
#include "conf2/browse.h"
#include "conf2/schema.h"

JNIEXPORT jlong JNICALL Java_org_conf2_schema_yang_ModuleLoader_loadModule
  (JNIEnv *env, jclass dloader_class, jobject j_streamsource, jstring resource,
   jobject j_module_browser) {
  void* stream_source_hnd_id = conf2_handle_new(j_streamsource, NULL);
  void *go_stream_source_hnd_id = conf2_new_driver_resource_source(stream_source_hnd_id,
    &conf2j_open_stream, &conf2j_read_stream);

  jobject j_g_module_browser = (*env)->NewGlobalRef(env, j_module_browser);
  void *module_browser_hnd_id = conf2_handle_new(j_g_module_browser, &conf2j_release_global_ref);

  const char *resourceStr = (*env)->GetStringUTFChars(env, resource, 0);
  void *module = conf2_load_module(
        &conf2j_browse_enter,
        &conf2j_browse_iterate,
        &conf2j_browse_read,
        &conf2j_browse_edit,
        &conf2j_browse_choose,
        &conf2j_browse_event,
        &conf2j_browse_find,
        &conf2j_browse_select_root,
        module_browser_hnd_id,
        go_stream_source_hnd_id,
        (char *)resourceStr);
  (*env)->ReleaseStringUTFChars(env, resource, resourceStr);
  return (jlong) module;
}