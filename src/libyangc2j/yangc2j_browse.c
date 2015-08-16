#include <string.h>
#include <stdbool.h>
#include "libyangc2j/libyangc2j.h"
#include "org_conf2_yang_driver_DriverLoader.h"
#include "yang/driver/yangc2_browse.h"
#include "driver.h"

bool checkBrowseError(JNIEnv *env, GoInterface *err) {
  if ((*env)->ExceptionCheck(env)) {
    jthrowable exception = (*env)->ExceptionOccurred(env);
    char *msg = get_exception_message(env, exception);
    *err = yangc2_new_browse_error(msg);
    (*env)->ExceptionClear(env);
    return true;
  }

  return false;
}

void* java_browse_root_selector(void *browser_handle, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();

  // get java browser instance
  jobject j_browser = browser_handle;

  // call rootSelector method
  jclass j_browser_cls = (*env)->FindClass(env, "org/conf2/yang/browse/Browser");
  if (checkBrowseError(env, err)) {
    return 0;
  }
  jmethodID selector_method = (*env)->GetMethodID(env, j_browser_cls, "getRootSelector", "()Lorg/conf2/yang/browse/Selection;");
  if (checkBrowseError(env, err)) {
    return 0;
  }

  jobject j_selector = (*env)->CallObjectMethod(env, j_browser, selector_method);
  if (checkBrowseError(env, err)) {
    return 0;
  }

  return j_selector
}

jmethodID get_adapter_method(JNIEnv *env, GoInterface *err, char *method_name, char *signature) {
  jclass adapter_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowserAdaptor");
  if (checkBrowseError(env, err)) {
    return 0;
  }
  jmethodID enter_method = (*env)->GetMethodID(env, j_browser_cls, method_name, signature);
  if (checkBrowseError(env, err)) {
    return 0;
  }
}

void* java_browse_enter(void *selection_handle, char *ident, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);

  jmethodID enter_method = get_adapter_method(env, err, "enter", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)V");
  if (enter_method == 0) {
    return 0
  }

  // null is valid
  jobject j_child_selector = (*env)->CallObjectMethod(env, enter_method, j_selection, j_ident);
  return j_child_selector;
}

short java_browse_iterate(void *selection_handle, char *ident, char *encodedKeys, short first, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  jobject j_encoded_keys = (*env)->NewStringUTF(env, encodedKeys);

  jmethodID iterate_method = get_adapter_method(env, err, "iterate", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;Ljava/lang/String;B)V");
  if (iterate_method == 0) {
    return 0
  }

  jboolean j_first = (jboolean) first;
  jboolean j_has_more = (*env)->CallBooleanMethod(env, iterate_method, j_selection, j_ident, j_encoded_keys, j_first);
  return (short)j_has_more;
}

void java_browse_read(void *selection_handle, char *ident, yangc2_browse_value *val, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);

  jmethodID read_method = get_adapter_method(env, err, "read", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)org/conf2/yang/browse/BrowseValue");
  if (read_method == 0) {
    return 0
  }

  jobject j_value = (*env)->CallObjectMethod(env, read_method, j_selection, j_ident);

  jclass value_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowseValue");
  if (checkBrowseError(env, err)) {
    return 0;
  }

  jfieldID s_val_field = (*env)->GetFieldID(env, value_cls, "valType", "org/conf2/yang/ValueType;");
  if (checkBrowseError(env, err)) {
    return 0;
  }

  val->val_type = (*env)->GetIntField(env, j_value, s_val_field);
  if (checkBrowseError(env, err)) {
    return 0;
  }

  switch val->val_type {
    case STRING:
      jfieldID s_val_field = (*env)->GetFieldID(env, value_cls, "str", "java/lang/String;");
      if (checkBrowseError(env, err)) {
        return 0;
      }
      jobject j_str = (*env)->GetObjectField(env, s_val_field);
      const char *c_str = (*env)->GetStringUTFChars(env, j_str, 0);
      strncpy(val->s, cResourceId, sizeof(val->s));
      (*env)->ReleaseStringUTFChars(env, j_str, c_str);

      // TODO: Support longer strings

      break;
    case INT32:
      jfieldID i_val_field = (*env)->GetFieldID(env, value_cls, "int32", "I");
      if (checkBrowseError(env, err)) {
        return 0;
      }
      jint j_i = (*env)->GetIntField(env, i_val_field);
      val->i = (int)j_i;
      break;
    case BOOLEAN:
      jfieldID b_val_field = (*env)->GetFieldID(env, value_cls, "bool", "B");
      if (checkBrowseError(env, err)) {
        return 0;
      }
      jboolean j_bool = (*env)->GetBooleanField(env, b_val_field);
      val->b = (short)j_bool;
      break;
    case EMPTY:
      break;
    default:
      *err = yangc2_new_browse_error("Unsupported type");
      break;
  }
}


void yangc2_browse_edit(void *selection_handle, char *ident, int op, yangc2_browse_value *val, void *browse_err) {
  GoInterface *err = (GoInterface *) browse_err;
  JNIEnv* env = getCurrentJniEnv();
  jobject j_selection = selection_handle;
  jobject j_ident = (*env)->NewStringUTF(env, ident);
  jobject j_value = NULL;
  jint j_op = (jint) op;

  jclass value_cls = (*env)->FindClass(env, "org/conf2/yang/browse/BrowseValue");
  if (checkBrowseError(env, err)) {
    return;
  }

  if val != NULL {
      switch val->val_type {
        case STRING: {
          jfieldID val_field = (*env)->GetMethodID(env, value_cls, "Str", "(Lorg/conf2/yang/browse/BrowseValue;)java/lang/String;");
          if (checkBrowseError(env, err)) {
            return;
          }
          jobject s_value = (*env)->NewStringUTF(env, val->s);
          jobject j_value = (*env)->CallObjectMethod(env, val_field, s_value);
        }
        case INT32: {
          jfieldID val_field = (*env)->GetMethodID(env, value_cls, "Int32", "(Lorg/conf2/yang/browse/BrowseValue;)I");
          if (checkBrowseError(env, err)) {
            return;
          }
          jobject j_value = (*env)->CallObjectMethod(env, val_field, val->i);
        }
        case BOOLEAN: {
          jfieldID val_field = (*env)->GetMethodID(env, value_cls, "Bool", "(Lorg/conf2/yang/browse/BrowseValue;)B");
          if (checkBrowseError(env, err)) {
            return;
          }
          jobject j_value = (*env)->CallObjectMethod(env, val_field, val->b);
        }
        default:
          *err = yangc2_new_browse_error("Unsupported type");
          return;
      }
  }

  jmethodID edit_method = get_adapter_method(env, err, "edit", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)org/conf2/yang/browse/BrowseValue");
  if (edit_method == 0) {
    return;
  }

  (*env)->CallVoidMethod(env, edit_method, j_selection, j_ident, j_op, j_value);
  if (checkBrowseError(env, err)) {
    return;
  }
}

void yangc2_browse_choose(void *selection_handle, char *ident, char *resolved_case, int resolved_case_max_len, void *browse_err) {
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_ident = (*env)->NewStringUTF(env, ident);

    jmethodID choose_method = get_adapter_method(env, err, "choose", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)Ljava/lang/String;");
    if (choose_method == 0) {
      return;
    }

    jobject j_choice = (*env)->CallObjectMethod(env, choose_method, j_selection, j_ident);
    if (checkBrowseError(env, err)) {
      return;
    }

    const char *c_choice = (*env)->GetStringUTFChars(env, j_choice, 0);
    strncpy(resolved_case, c_choice, resolved_case_max_len);
    (*env)->ReleaseStringUTFChars(env, j_choice, c_choice);
}

void yangc2_browse_exit(void *selection_handle, char *ident, void *browse_err) {
    GoInterface *err = (GoInterface *) browse_err;
    JNIEnv* env = getCurrentJniEnv();
    jobject j_selection = selection_handle;
    jobject j_ident = (*env)->NewStringUTF(env, ident);

    jmethodID exit_method = get_adapter_method(env, err, "exit", "(Lorg/conf2/yang/browse/Selection;Ljava/lang/String;)V");
    if (exit_method == 0) {
      return;
    }

    (*env)->CallVoidMethod(env, exit_method, j_selection, j_ident);
    if (checkBrowseError(env, err)) {
      return;
    }
}



JNIEXPORT jobject JNICALL Java_org_conf2_yang_driver_DriverLoader_loadModule
  (JNIEnv *env, jclass dloader_class, jobject datasource_hnd, jstring resource, jobject module_browser) {

  GoInterface browser = yangc2_browse_new_yang_browser(
        java_browse_enter,
        java_browse_iterate,
        java_browse_read,
        java_brose_edit,
        java_browse_choose,
        java_browse_exit,
        java_browse_root_selector,
        module_browser);

  // TODO: Determine if GC could collect browser instance during process?
  jobject browser_hnd = makeDriverHandle(env, browser);
  return browser_hnd;
}
