#include <stdio.h>
#include "yang-c2/browse.h"

typedef struct { void *data; } GoSlice;

// Bridge functions to call C function pointer in a given language data browsers
void *yangc2_browse_root_selector(yangc2_browse_root_selector_impl impl_func, void *browser_handle, void *browse_err) {
    return (*impl_func)(browser_handle, browse_err);
}

void *yangc2_browse_enter(yangc2_browse_enter_impl impl_func, void *selection_handle, char *ident, short *found, void *browse_err) {
    return (*impl_func)(selection_handle, ident, found, browse_err);
}

short yangc2_browse_iterate(yangc2_browse_iterate_impl impl_func, void *selection_handle, char *encodedKeys, short first, void *browse_err) {
    return (*impl_func)(selection_handle, encodedKeys, first, browse_err);
}

void yangc2_browse_read(yangc2_browse_read_impl impl_func, void *selection_handle, char *ident, struct yangc2_value *val, void *browse_err) {
    return (*impl_func)(selection_handle, ident, val, browse_err);
}

void yangc2_browse_edit(yangc2_browse_edit_impl impl_func, void *selection_handle, char *ident, int op, struct yangc2_value *val, void *browse_err) {
    return (*impl_func)(selection_handle, ident, op, val, browse_err);
}

char *yangc2_browse_choose(yangc2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err) {
    return (*impl_func)(selection_handle, ident, browse_err);
}

void yangc2_browse_exit(yangc2_browse_exit_impl impl_func, void *selection_handle, char *ident, void *browse_err) {
    return (*impl_func)(selection_handle, ident, browse_err);
}