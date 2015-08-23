#include <string.h>
#include <stdbool.h>
#include <stdlib.h>
#include "libyangc2j/java_yang.h"
#include "org_conf2_yang_driver_DriverLoader.h"
#include "yang/driver/yangc2_browse.h"
#include "yang/driver.h"

void* java_browse_root_selector(void *browser_handle, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();

  // get java browser instance
  jobject j_browser = browser_handle;

  // call rootSelector method
  jclass j_browser_cls = (*env)->FindClass(env, "org/conf2/yang/browse/Browser");
  if (checkDriverError(env, err)) {
    return NULL;
  }
  jmethodID selector_method = (*env)->GetMethodID(env, j_browser_cls, "getRootSelector", "()Lorg/conf2/yang/browse/Selection;");
  if (checkDriverError(env, err)) {
    return NULL;
  }

  jobject j_selector = (*env)->CallObjectMethod(env, j_browser, selector_method);
  if (checkDriverError(env, err)) {
    return NULL;
  }

  return j_selector;
}

//typedef struct _java_string_release {
//  JEnv *env;
//  jobject j_str;
//  const char *c_str;
//} java_string_release;
//
//void java_release_string(void *data) {
//   java_string_release *release = (java_string_release *)data;
//  (*(release->env))->ReleaseStringUTFChars(release->env, release->j_str, release->c_str);
//  free(release);
//}

yangc2j_method get_adapter_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  yangc2j_method m;
  m.cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowserAdaptor");
  if (checkDriverError(env, err)) {
    return;
  }
  m.methodId = (*env)->GetMethodID(env, m.cls, method_name, signature);
  if (checkDriverError(env, err)) {
    return;
  }
  return m;
}

void* java_browse_enter(void *selection_handle, char *ident, short *found, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);

  yangc2j_method enter = get_adapter_method(env, err, "enter", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)V");
  if (enter.methodId == NULL) {
    return NULL;
  }

  // null is valid
  jobject j_child_selector = (*env)->CallStaticObjectMethod(env, enter.cls, enter.methodId,
    j_selection, j_ident);

  if (j_child_selector != NULL) {
    jclass selection_cls = (*env)->FindClass(env, "org/conf2/yang/browse/Selection");
    if (checkDriverError(env, err)) {
      return;
    }

    jfieldID found_field = (*env)->GetFieldID(env, selection_cls, "found", "B");
    if (checkDriverError(env, err)) {
      return;
    }
    jboolean j_found = (*env)->GetBooleanField(env, j_selection, found_field);
    *found = (short)j_found;
  }

  return j_child_selector;
}

short java_browse_iterate(void *selection_handle, char *ident, char *encodedKeys, short first, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  jobject j_encoded_keys = (*env)->NewStringUTF(env, encodedKeys);

  yangc2j_method iterate = get_adapter_method(env, err, "iterate", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;Ljava/lang/String;B)V");
  if (iterate.methodId == NULL) {
    return false;
  }

  jboolean j_first = (jboolean) first;
  jboolean j_has_more = (*env)->CallStaticBooleanMethod(env, iterate.cls, iterate.methodId, j_selection, j_ident, j_encoded_keys, j_first);
  return (short)j_has_more;
}

void java_browse_read(void *selection_handle, char *ident, struct yangc2_browse_value *val, void *browse_err) {
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

  switch (val->val_type) {
    case STRING: {
      jfieldID s_val_field = (*env)->GetFieldID(env, value_cls, "str", "java/lang/String;");
      if (checkDriverError(env, err)) {
        return;
      }
      jobject j_str = (*env)->GetObjectField(env, j_value, s_val_field);
      const char* str = (*env)->GetStringUTFChars(env, j_str, 0);
      val->data = malloc((size_t)strlen(str));
      val->str = (char *)val->data;
      strcpy(val->str, str);
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
      jfieldID b_val_field = (*env)->GetFieldID(env, value_cls, "bool", "B");
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

void java_browse_edit(void *selection_handle, char *ident, int op, struct yangc2_browse_value *val, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  jobject j_value = NULL;
  jint j_op = (jint) op;

  jclass value_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowseValue");
  if (checkDriverError(env, err)) {
    return;
  }

  if (val != NULL) {
      switch (val->val_type) {
        case STRING: {
          jfieldID val_field = (*env)->GetFieldID(env, value_cls, "Str", "(Lorg/conf2/yang/browse/BrowseValue;)java/lang/String;");
          if (checkDriverError(env, err)) {
            return;
          }
          jobject s_value = (*env)->NewStringUTF(env, val->str);
          (*env)->SetObjectField(env, j_value, val_field, s_value);
        }
        case INT32: {
          jfieldID val_field = (*env)->GetFieldID(env, value_cls, "Int32", "(Lorg/conf2/yang/browse/BrowseValue;)I");
          if (checkDriverError(env, err)) {
            return;
          }
          (*env)->SetIntField(env, j_value, val_field, val->int32);
        }
        case BOOLEAN: {
          jfieldID val_field = (*env)->GetFieldID(env, value_cls, "Bool", "(Lorg/conf2/yang/browse/BrowseValue;)B");
          if (checkDriverError(env, err)) {
            return;
          }
          (*env)->SetBooleanField(env, j_value, val_field, val->boolean);
        }
        default:
          *err = yangc2_new_driver_error("Unsupported type");
          return;
      }
  }

  yangc2j_method edit = get_adapter_method(env, err, "edit", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)org/conf2/yang/browse/BrowseValue");
  if (edit.methodId == NULL) {
    return;
  }

  (*env)->CallStaticVoidMethod(env, edit.cls, edit.methodId, j_selection, j_ident, j_op, j_value);
  if (checkDriverError(env, err)) {
    return;
  }
}

char *java_browse_choose(void *selection_handle, char *ident, void *browse_err) {
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

JNIEXPORT jobject JNICALL Java_org_conf2_yang_driver_DriverLoader_loadModule
  (JNIEnv *env, jclass dloader_class, jobject datasource_hnd, jstring resource, jobject module_browser) {

  GoInterface rs_source;
printf("C loadModule - before resolve driver handle, %p\n", datasource_hnd);
  resolveDriverHandle(env, datasource_hnd, &rs_source);
  const char *resourceStr = (*env)->GetStringUTFChars(env, resource, 0);
printf("C loadModule - after resolve driver\n");
  GoInterface module = yangc2_load_module(
        &java_browse_enter,
        &java_browse_iterate,
        &java_browse_read,
        &java_browse_edit,
        &java_browse_choose,
        &java_browse_exit,
        &java_browse_root_selector,
        module_browser,
        rs_source,
        (char *)resourceStr);
printf("C loadModule - after yangc2_load_module\n");
  (*env)->ReleaseStringUTFChars(env, resource, resourceStr);
  if (module.v != NULL) {
    // throw proper error
    return NULL;
  }

  // TODO: Determine if GC could collect browser instance during process?
  jobject module_hnd = makeDriverHandle(env, module);
  return module_hnd;
}
