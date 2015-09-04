#ifndef YANGC2_BROWSE_H
#define YANGC2_BROWSE_H

struct GoSlice;

// matches list in meta.go
enum yangc2_format {
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
    FMT_UNION
};

struct yangc2_value {
    enum yangc2_format format;
    // PERFORMANCE TODO: See if you can use a union.  Unclear of Go's integration level.
    int int32;
    short is_list;
    int list_len;
    short boolean;
    char  *cstr;
    void *data;
    int data_len;
    void *handle;
};

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to convert data streams.  You will pass your function pointer when calling
// yangc2_browse_new_browser
typedef void* (*yangc2_browse_root_selector_impl)(void *browser_handle, void *browse_err);
typedef void* (*yangc2_browse_enter_impl)(void *selection_handle, char *ident, short *found, void *browse_err);
typedef short (*yangc2_browse_iterate_impl)(void *selection_handle, char *keys, short first, void *browse_err);
typedef void (*yangc2_browse_read_impl)(void *selection_handle, char *ident, struct yangc2_value* val, void *browse_err);
typedef void (*yangc2_browse_edit_impl)(void *selection_handle, char *ident, int op, struct yangc2_value* val, void *browse_err);
typedef char* (*yangc2_browse_choose_impl)(void *selection_handle, char *ident, void *browse_err);
typedef void (*yangc2_browse_exit_impl)(void *selection_handle, char *ident, void *browse_err);

#endif
