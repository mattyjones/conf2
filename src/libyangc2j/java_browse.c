#include <string.h>
#include <stdbool.h>
#include <stdlib.h>
#include "yang-c2/java.h"
#include "org_conf2_yang_driver_DriverLoader.h"
#include "yang-c2/browse.h"
#include "yang-c2/driver.h"

yangc2j_method get_adapter_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  yangc2j_method m;
  m.methodId = NULL;
  m.cls = NULL;
  m.cls = (*env)->FindClass(env, "org/conf2/yang/driver/BrowserAdaptor");
  if (checkDriverError(env, err)) {
    return m;
  }
  m.methodId = (*env)->GetStaticMethodID(env, m.cls, method_name, signature);
  if (checkDriverError(env, err)) {
    return m;
  }
  return m;
}

void* java_browse_root_selector(void *browser_handle, void *browse_err) {
printf("java_browse.c:root HERE\n");
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();

  // get java browser instance
  jobject j_browser = browser_handle;

  yangc2j_method root_selector = get_adapter_method(env, err, "getRootSelector",
    "(Lorg/conf2/yang/browse/Browser;)Lorg/conf2/yang/browse/Selection;");
  if (root_selector.methodId == NULL) {
printf("java_browse.c:root ERR\n");
    return NULL;
  }

  jobject j_selector = (*env)->CallStaticObjectMethod(env, root_selector.cls, root_selector.methodId, j_browser);
  jobject j_g_selector = (*env)->NewGlobalRef(env, j_selector);
  void *root_selector_hnd_id = yangc2_handle_new(j_g_selector, &java_release_global_ref);
printf("java_browse.c:root j_g_selector=%p, root_selector_hnd_id=%p\n", j_g_selector, root_selector_hnd_id);
  return root_selector_hnd_id;
}

void* java_browse_enter(void *selection_handle, char *ident, short *found, void *browse_err) {
printf("java_browse.c:enter selection_handle=%p, ident=%s\n", selection_handle, ident);
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  void *child_selector_hnd_id = NULL;

  yangc2j_method enter = get_adapter_method(env, err, "enter", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)Lorg/conf2/yang/browse/Selection;");
  if (enter.methodId == NULL) {
    return NULL;
  }

  // null is valid
  jobject j_child_selector = (*env)->CallStaticObjectMethod(env, enter.cls, enter.methodId,
    j_selection, j_ident);

  if (j_child_selector != NULL) {
    jobject j_g_child_selector = (*env)->NewGlobalRef(env, j_child_selector);
    child_selector_hnd_id = yangc2_handle_new(j_g_child_selector, java_release_global_ref);
  }

  jclass selection_cls = (*env)->FindClass(env, "org/conf2/yang/browse/Selection");
  if (checkDriverError(env, err)) {
    return;
  }

  jfieldID found_field = (*env)->GetFieldID(env, selection_cls, "found", "Z");
  if (checkDriverError(env, err)) {
    return;
  }
  jboolean j_found = (*env)->GetBooleanField(env, j_selection, found_field);
  *found = (short)j_found;
printf("java_browse.c:enter FOUND=%d\n", j_found);

  return child_selector_hnd_id;
}

short java_browse_iterate(void *selection_handle, char *encodedKeys, short first, void *browse_err) {
printf("java_browse.c:iterate selection_handle=%p\n", selection_handle);
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_encoded_keys = (*env)->NewStringUTF(env, encodedKeys);

  yangc2j_method iterate = get_adapter_method(env, err, "iterate", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;Z)Z");
  if (iterate.methodId == NULL) {
    return false;
  }

  jboolean j_first = (jboolean) first;
  jboolean j_has_more = (*env)->CallStaticBooleanMethod(env, iterate.cls, iterate.methodId, j_selection, j_encoded_keys, j_first);
  if (checkDriverError(env, err)) {
    return;
  }

  return (short)j_has_more;
}

void java_browse_read(void *selection_handle, char *ident, struct yangc2_browse_value *val, void *browse_err) {
printf("java_browse.c:read selection_handle=%p, ident=%s\n", selection_handle, ident);
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);

  yangc2j_method read = get_adapter_method(env, err, "read", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)org/conf2/yang/browse/BrowseValue");
  if (read.methodId == NULL) {
    return;
  }

  jobject j_value = (*env)->CallStaticObjectMethod(env, read.cls, read.methodId, j_selection, j_ident);

  jclass value_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowseValue");
  if (checkDriverError(env, err)) {
    return;
  }

  jfieldID s_val_field = (*env)->GetFieldID(env, value_cls, "valType", "org/conf2/yang/ValueType;");
  if (checkDriverError(env, err)) {
    return;
  }

  val->val_type = (*env)->GetIntField(env, j_value, s_val_field);
  if (checkDriverError(env, err)) {
    return;
  }

  // TODO: LeafList

  switch (val->val_type) {
    case STRING: {
      jfieldID s_val_field = (*env)->GetFieldID(env, value_cls, "str", "java/lang/String;");
      if (checkDriverError(env, err)) {
        return;
      }
      jobject j_str = (*env)->GetObjectField(env, j_value, s_val_field);
      if (j_str != NULL) {
          java_string_chars* chars = java_new_string_chars(j_str);
          val->handle = chars->handle;
          val->str = (char *)chars->chars;
      }
      break;
    }
    case INT32: {
      jfieldID i_val_field = (*env)->GetFieldID(env, value_cls, "int32", "I");
      if (checkDriverError(env, err)) {
        return;
      }
      jint j_i = (*env)->GetIntField(env, j_value, i_val_field);
      val->int32 = (int)j_i;
      break;
    }
    case BOOLEAN: {
      jfieldID b_val_field = (*env)->GetFieldID(env, value_cls, "bool", "Z");
      if (checkDriverError(env, err)) {
        return;
      }
      jboolean j_bool = (*env)->GetBooleanField(env, j_value, b_val_field);
      val->boolean = (short)j_bool;
      break;
    }
    case EMPTY:
      break;
    default: {
      *err = yangc2_new_driver_error("Unsupported type");
      break;
    }
  }
}

jobject java_value(JNIEnv *env, struct yangc2_browse_value *val, GoInterface *err) {
  jobject j_value = NULL;
  jclass value_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowseValue");
  if (checkDriverError(env, err)) {
    return;
  }
  if (val->islist) {
      switch (val->val_type) {
        case STRING: {
            jclass str_class = (*env)->FindClass(env, "java/lang/String");
            jobjectArray j_strlist = (*env)->NewObjectArray(env, val->listlen, str_class, NULL);
            if (j_strlist == NULL) {
                *err = yangc2_new_driver_error("Out of memory");
                return NULL;
            }
            int i;
            for (i = 0; i < val->listlen; i++) {
                (*env)->SetObjectArrayElement(env, j_strlist, i, (*env)->NewStringUTF(env, val->strlist[i]));
            }

            jmethodID factory_method = (*env)->GetStaticMethodID(env, value_cls, "Strlist", "([Ljava/lang/String;)Lorg/conf2/yang/browse/BrowseValue;");
            if (checkDriverError(env, err)) {
                return;
            }
            j_value = (*env)->CallStaticObjectMethod(env, value_cls, factory_method, j_strlist);
            if (checkDriverError(env, err)) {
                return;
            }
            break;
        }
        case INT32:
            // TODO
            break;
        case BOOLEAN:
            // TODO
            break;
      }
  } else {
      switch (val->val_type) {
        case STRING: {
          jobject s_value = (*env)->NewStringUTF(env, val->str);

          jmethodID factoryMethod = (*env)->GetStaticMethodID(env, value_cls, "Str", "(Ljava/lang/String;)Lorg/conf2/yang/browse/BrowseValue;");
          if (checkDriverError(env, err)) {
            return;
          }

          j_value = (*env)->CallStaticObjectMethod(env, value_cls, factoryMethod, s_value);
          if (checkDriverError(env, err)) {
            return;
          }
          break;
        }
        case INT32: {
          jmethodID factoryMethod = (*env)->GetStaticMethodID(env, value_cls, "Int32", "(I)Lorg/conf2/yang/browse/BrowseValue;");
          if (checkDriverError(env, err)) {
            return;
          }

          j_value = (*env)->CallStaticObjectMethod(env, value_cls, factoryMethod, val->int32);
          if (checkDriverError(env, err)) {
            return;
          }
          break;
        }
        case BOOLEAN: {
          jmethodID factoryMethod = (*env)->GetStaticMethodID(env, value_cls, "Bool", "(Z)Lorg/conf2/yang/browse/BrowseValue;");
          if (checkDriverError(env, err)) {
            return;
          }

          j_value = (*env)->CallStaticObjectMethod(env, value_cls, factoryMethod, (jboolean)val->boolean);
          if (checkDriverError(env, err)) {
            return;
          }
          break;
        }
        case EMPTY:
          break;
        default:
          *err = yangc2_new_driver_error("Unsupported type");
          return;
      }
  }
  return j_value;
}

void java_browse_edit(void *selection_handle, char *ident, int op, struct yangc2_browse_value *val, void *browse_err) {
printf("java_browse.c:edit selection_handle=%p, op=%d, ident=%s\n", selection_handle, op, ident);
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = NULL;
  if (ident != NULL) {
    j_ident = (*env)->NewStringUTF(env, ident);
  }
  jobject j_value = NULL;
  jint j_op = (jint) op;

  if (val != NULL) {
    j_value = java_value(env, val, err);
  }

  yangc2j_method edit = get_adapter_method(env, err, "edit", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;ILorg/conf2/yang/browse/BrowseValue;)V");
  if (edit.methodId == NULL) {
    return;
  }

printf("java_browse.c:edit CallStaticVoidMethod j_selection=%p, j_value=%p\n", j_selection, j_value);
  (*env)->CallStaticVoidMethod(env, edit.cls, edit.methodId, j_selection, j_ident, j_op, j_value);
  if (checkDriverError(env, err)) {
printf("java_browse.c:edit ERROR CallStaticVoidMethod\n");
    return;
  }
printf("java_browse.c:edit LEAVING\n");
}

char *java_browse_choose(void *selection_handle, char *ident, void *browse_err) {
printf("java_browse.c:choose selection_handle=%p, ident=%s\n", selection_handle, ident);
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_ident = (*env)->NewStringUTF(env, ident);

    yangc2j_method choose = get_adapter_method(env, err, "choose", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)Ljava/lang/String;");
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

void java_browse_exit(void *selection_handle, char *ident, void *browse_err) {
printf("java_browse.c:exit selection_handle=%p, ident=%s\n", selection_handle, ident);
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_ident = (*env)->NewStringUTF(env, ident);

    yangc2j_method exit = get_adapter_method(env, err, "exit", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)V");
    if (exit.methodId == 0) {
      return;
    }

    (*env)->CallStaticVoidMethod(env, exit.cls, exit.methodId, j_selection, j_ident);
    if (checkDriverError(env, err)) {
      return;
    }
}

JNIEXPORT jlong JNICALL Java_org_conf2_yang_driver_Driver_loadModule
  (JNIEnv *env, jclass dloader_class, jobject j_streamsource, jstring resource,
   jobject j_module_browser) {

  void* stream_source_hnd_id = yangc2_handle_new(j_streamsource, NULL);
printf("java_stream.c:newStreamSource stream_source_hnd_id=%p\n", stream_source_hnd_id);
  void *go_stream_source_hnd_id = yangc2_new_driver_resource_source(stream_source_hnd_id, &java_open_stream,
    &java_read_stream);

  jobject j_g_module_browser = (*env)->NewGlobalRef(env, j_module_browser);
  void *module_browser_hnd_id = yangc2_handle_new(j_g_module_browser, &java_release_global_ref);

  const char *resourceStr = (*env)->GetStringUTFChars(env, resource, 0);
  void *module = yangc2_load_module(
        &java_browse_enter,
        &java_browse_iterate,
        &java_browse_read,
        &java_browse_edit,
        &java_browse_choose,
        &java_browse_exit,
        &java_browse_root_selector,
        module_browser_hnd_id,
        go_stream_source_hnd_id,
        (char *)resourceStr);
  (*env)->ReleaseStringUTFChars(env, resource, resourceStr);
  return (jlong) module;
}
