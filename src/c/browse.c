#include <stdio.h>
#include "conf2/browse.h"

typedef struct { void *data; } GoSlice;

// Bridge functions to call C function pointer in a given language data browsers
void *conf2_browse_selector(conf2_browse_selector_impl impl_func, void *browser_handle, char *path, void *browse_err) {
    return (*impl_func)(browser_handle, path, browse_err);
}

void *conf2_browse_enter(conf2_browse_enter_impl impl_func, void *selection_handle, char *ident, short create, void *browse_err) {
    return (*impl_func)(selection_handle, ident, create, browse_err);
}

void *conf2_browse_iterate(conf2_browse_iterate_impl impl_func, void *selection_handle, short create, void *key_data, int key_data_len, short first, void *browse_err) {
    return (*impl_func)(selection_handle, create, key_data, key_data_len, first, browse_err);
}

void *conf2_browse_read(conf2_browse_read_impl impl_func, void *selection_handle, char *ident, void **val_data_ptr, int* val_data_len_ptr, void *browse_err) {
    return (*impl_func)(selection_handle, ident, val_data_ptr, val_data_len_ptr, browse_err);
}

void conf2_browse_edit(conf2_browse_edit_impl impl_func, void *selection_handle, char *ident, void *val_data, int val_data_len, void *browse_err) {
    return (*impl_func)(selection_handle, ident, val_data, val_data_len, browse_err);
}

char *conf2_browse_choose(conf2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err) {
    return (*impl_func)(selection_handle, ident, browse_err);
}

void conf2_browse_event(conf2_browse_event_impl impl_func, void *selection_handle, int eventId, void *browse_err) {
    return (*impl_func)(selection_handle, eventId, browse_err);
}
