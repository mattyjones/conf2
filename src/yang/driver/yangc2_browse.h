#ifndef YANGC2_BROWSE_H
#define YANGC2_BROWSE_H

enum yangc2_browse_value_type {
    EMPTY,
    BINARY,
    BITS,
    BOOLEAN,
    DECIMAL64,
    ENUMERATION,
    IDENTITYDEF,
    INSTANCE_IDENTIFIER,
    INT8,
    INT16,
    INT32,
    INT64,
    LEAFREF,
    STRING,
    UINT8,
    UINT16,
    UINT32,
    UINT64,
    UNION
};

struct yangc2_browse_value {
    enum yangc2_browse_value_type val_type;
    // PERFORMANCE TODO: See if you can use a union.  Unclear of Go's integration level.
    int int32;
    short boolean;
    char  *str;
    void *data;
};

// Each language (python, java, etc.) will implement at most one of each of these functions
// pointers to convert data streams.  You will pass your function pointer when calling
// yangc2_browse_new_browser
typedef void* (*yangc2_browse_enter_impl)(void *selection_handle, char *ident, short *found, void *browse_err);
typedef short (*yangc2_browse_iterate_impl)(void *selection_handle, char *ident, char *keys, short first, void *browse_err);
typedef void (*yangc2_browse_read_impl)(void *selection_handle, char *ident, struct yangc2_browse_value* val, void *browse_err);
typedef void (*yangc2_browse_edit_impl)(void *selection_handle, char *ident, int op, struct yangc2_browse_value* val, void *browse_err);
typedef char* (*yangc2_browse_choose_impl)(void *selection_handle, char *ident, void *browse_err);
typedef void (*yangc2_browse_exit_impl)(void *selection_handle, char *ident, void *browse_err);
typedef void* (*yangc2_browse_root_selector_impl)(void *browser_handle, void *browse_err);



#endif
