#include <string.h>
#include <stdbool.h>
#include <stdlib.h>
#include "conf2/java.h"
#include "conf2/browse.h"
#include "conf2/schema.h"

conf2j_method conf2j_static_adapter_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  return conf2j_static_method(env, err, "org/conf2/schema/driver/BrowserAdaptor", method_name, signature);
}

void* conf2j_browse_select_root(void *browser_handle, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();

  // get java browser instance
  jobject j_browser = browser_handle;

  conf2j_method selector = conf2j_static_adapter_method(env, err, "getSelector",
    "(Lorg/conf2/schema/browse/Browser;)Lorg/conf2/schema/browse/Selection;");
  if (selector.methodId == NULL) {
    return NULL;
  }

  jobject j_selector = (*env)->CallStaticObjectMethod(env, selector.cls, selector.methodId, j_browser);
  jobject j_g_selector = (*env)->NewGlobalRef(env, j_selector);
  void *selector_hnd_id = conf2_handle_new(j_g_selector, &conf2j_release_global_ref);
  return selector_hnd_id;
}

void *conf2j_browse_enter(void *selection_handle, char *ident, short create, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  void *child_selector_hnd_id = NULL;

  conf2j_method enter = conf2j_static_adapter_method(env, err, "enter", "(Lorg/conf2/schema/browse/Selection;Ljava/lang/String;Z)Lorg/conf2/schema/browse/Node;");
  if (enter.methodId == NULL) {
    return NULL;
  }

  // null is valid
  jobject j_child_selector = (*env)->CallStaticObjectMethod(env, enter.cls, enter.methodId,
    j_selection, j_ident, create > 0);

  if (j_child_selector != NULL) {
    jobject j_g_child_selector = (*env)->NewGlobalRef(env, j_child_selector);
    child_selector_hnd_id = conf2_handle_new(j_g_child_selector, conf2j_release_global_ref);
  }

  return child_selector_hnd_id;
}

void *conf2j_browse_read(void *selection_handle, char *ident, void **val_data_ptr, int *val_data_len, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  void *handle = NULL;

  conf2j_method read = conf2j_static_adapter_method(env, err, "read", "(Lorg/conf2/schema/browse/Selection;Ljava/lang/String;)Ljava/nio/ByteBuffer;");
  if (read.methodId == NULL) {
    return;
  }

  jobject j_value_buff = (*env)->CallStaticObjectMethod(env, read.cls, read.methodId, j_selection, j_ident);
  if (checkDriverError(env, err)) {
    return;
  }

  if (j_value_buff == NULL) {
    *val_data_ptr = NULL;
    *val_data_len = 0;
  } else {
    jobject j_g_value_buff = (*env)->NewGlobalRef(env, j_value_buff);
    *val_data_ptr = (*env)->GetDirectBufferAddress(env, j_g_value_buff);
    *val_data_len = (*env)->GetDirectBufferCapacity(env, j_g_value_buff);
    handle = conf2_handle_new(j_g_value_buff, &conf2j_release_global_ref);
  }
  return handle;
}

void *conf2j_browse_iterate(void *selection_handle, short create, void *key_data, int key_data_len, short first, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_encoded_key = NULL;
  void *child_selector_hnd_id = NULL;

  if (key_data != NULL && key_data_len > 0) {
      jobject j_encoded_key = (*env)->NewDirectByteBuffer(env, key_data, key_data_len);
  }

  conf2j_method iterate = conf2j_static_adapter_method(env, err, "iterate",
        "(Lorg/conf2/schema/browse/Selection;Ljava/nio/ByteBuffer;ZZ)Lorg/conf2/schema/browse/Node;");
  if (iterate.methodId == NULL) {
    return false;
  }

  jobject j_child_selector = (*env)->CallStaticObjectMethod(env, iterate.cls, iterate.methodId, j_selection, j_encoded_key,
    create > 0, first > 0);
  if (checkDriverError(env, err)) {
    return;
  }

  if (j_child_selector != NULL) {
    jobject j_g_child_selector = (*env)->NewGlobalRef(env, j_child_selector);
    child_selector_hnd_id = conf2_handle_new(j_g_child_selector, conf2j_release_global_ref);
  }

  return child_selector_hnd_id;
}

void conf2j_browse_edit(void *selection_handle, char *ident, void *val_data, int val_data_len, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = NULL;
  if (ident != NULL) {
    j_ident = (*env)->NewStringUTF(env, ident);
  }
  jobject j_encoded_value = NULL;

  if (val_data != NULL && val_data_len > 0) {
      jobject j_encoded_value = (*env)->NewDirectByteBuffer(env, val_data, val_data_len);
      if (checkDriverError(env, err)) {
        return;
      }
  }

  conf2j_method edit = conf2j_static_adapter_method(env, err, "edit", "(Lorg/conf2/schema/browse/Selection;Ljava/lang/String;Ljava/nio/ByteBuffer;)V");
  if (edit.methodId == NULL) {
    return;
  }

  (*env)->CallStaticVoidMethod(env, edit.cls, edit.methodId, j_selection, j_ident, j_encoded_value);
  if (checkDriverError(env, err)) {
    return;
  }
}

char *conf2j_browse_choose(void *selection_handle, char *ident, void *browse_err) {
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_ident = (*env)->NewStringUTF(env, ident);

    conf2j_method choose = conf2j_static_adapter_method(env, err, "choose", "(Lorg/conf2/schema/browse/Selection;Ljava/lang/String;)Ljava/lang/String;");
    if (choose.methodId == NULL) {
      return;
    }

    jobject j_choice = (*env)->CallStaticObjectMethod(env, choose.cls, choose.methodId, j_selection, j_ident);
    if (checkDriverError(env, err)) {
      return;
    }

    const char *c_choice = (*env)->GetStringUTFChars(env, j_choice, 0);
    char *resolved = (char *)malloc(strlen(c_choice));
    strcpy(resolved, c_choice);
    (*env)->ReleaseStringUTFChars(env, j_choice, c_choice);
    return resolved;
}

void conf2j_browse_event(void *selection_handle, int eventId, void *browse_err) {
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jint j_eventId = eventId;

    conf2j_method event = conf2j_static_adapter_method(env, err, "event", "(Lorg/conf2/schema/browse/Selection;I)V");
    if (event.methodId == 0) {
      return;
    }

    (*env)->CallStaticVoidMethod(env, event.cls, event.methodId, j_selection, j_eventId);
    if (checkDriverError(env, err)) {
      return;
    }
}

void conf2j_browse_find(void *selection_handle, char *path, void *browse_err) {
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_path = (*env)->NewStringUTF(env, path);

    conf2j_method event = conf2j_static_adapter_method(env, err, "find", "(Lorg/conf2/schema/browse/Selection;java/lang/String;)V");
    if (event.methodId == 0) {
      return;
    }

    (*env)->CallStaticVoidMethod(env, event.cls, event.methodId, j_selection, j_path);
    if (checkDriverError(env, err)) {
      return;
    }
}

void *conf2j_browse_new(JNIEnv *env, jlong module_hnd_id, jobject j_browser) {
  jobject j_g_browser = (*env)->NewGlobalRef(env, j_browser);
  void *browser_hnd_id = conf2_handle_new(j_g_browser, &conf2j_release_global_ref);
  return conf2_new_browser(
        &conf2j_browse_enter,
        &conf2j_browse_iterate,
        &conf2j_browse_read,
        &conf2j_browse_edit,
        &conf2j_browse_choose,
        &conf2j_browse_event,
        &conf2j_browse_find,
        &conf2j_browse_select_root,
        (void *)module_hnd_id,
        browser_hnd_id);
}


