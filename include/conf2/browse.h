#ifndef CONF2_BROWSE_H
#define CONF2_BROWSE_H

struct GoSlice;

// matches list in format.go
enum conf2_format {
    FMT_EMPTY,
    FMT_BINARY,
    FMT_BITS,
    FMT_BOOLEAN,
    FMT_DECIMAL64,
    FMT_ENUMERATION,
    FMT_IDENTITYDEF,
    FMT_INSTANCE_IDENTIFIER,
    FMT_INT8,
    FMT_INT16,
    FMT_INT32,
    FMT_INT64,
    FMT_LEAFREF,
    FMT_STRING,
    FMT_UINT8,
    FMT_UINT16,
    FMT_UINT32,
    FMT_UINT64,
    FMT_UNION,

    FMT_BINARY_LIST = 1025,
    FMT_BITS_LIST,
    FMT_BOOLEAN_LIST,
    FMT_DECIMAL64_LIST,
    FMT_ENUMERATION_LIST,
    FMT_IDENTITYDEF_LIST,
    FMT_INSTANCE_IDENTIFIER_LIST,
    FMT_INT8_LIST,
    FMT_INT16_LIST,
    FMT_INT32_LIST,
    FMT_INT64_LIST,
    FMT_LEAFREF_LIST,
    FMT_STRING_LIST,
    FMT_UINT8_LIST,
    FMT_UINT16_LIST,
    FMT_UINT32_LIST,
    FMT_UINT64_LIST,
    FMT_UNION_LIST,
};

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to convert data streams.  You will pass your function pointer when calling
// conf2_browse_new_browser
typedef void* (*conf2_browse_select_root_impl)(void *browser_handle, void *browse_err);
typedef void* (*conf2_browse_enter_impl)(void *selection_handle, char *ident, short create, void *browse_err);
typedef void* (*conf2_browse_iterate_impl)(void *selection_handle, short create, void *key_data, int key_data_len, short first, void *browse_err);
typedef void* (*conf2_browse_read_impl)(void *selection_handle, char *ident, void **val_data, int* val_data_len, void *browse_err);
typedef void (*conf2_browse_edit_impl)(void *selection_handle, char *ident, void *val_data, int val_data_len, void *browse_err);
typedef char* (*conf2_browse_choose_impl)(void *selection_handle, char *ident, void *browse_err);
typedef void (*conf2_browse_event_impl)(void *selection_handle, int eventId, void *browse_err);
typedef void (*conf2_browse_find_impl)(void *selection_handle, char *path, void *browse_err);
#endif
