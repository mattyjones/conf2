package driver

// #include "yang/driver/yangc2_browse.h"
// extern void* yangc2_browse_root_selector(yangc2_browse_root_selector_impl impl_func, void *browser_handle, void *browse_err);
// extern void* yangc2_browse_enter(yangc2_browse_enter_impl impl_func, void *selection_handle, char *ident, short *found, void *browse_err);
// extern short yangc2_browse_iterate(yangc2_browse_iterate_impl impl_func, void *selection_handle, char *ident, char *encodedKeys, short first, void *browse_err);
// extern void yangc2_browse_read(yangc2_browse_read_impl impl_func, void *selection_handle, char *ident, struct yangc2_browse_value *val, void *browse_err);
// extern void yangc2_browse_edit(yangc2_browse_edit_impl impl_func, void *selection_handle, char *ident, int op, struct yangc2_browse_value *val, void *browse_err);
// extern char *yangc2_browse_choose(yangc2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern void yangc2_browse_exit(yangc2_browse_exit_impl impl_func, void *selection_handle, char *ident, void *browse_err);
import "C"

import (
	"yang"
	"unsafe"
	"strings"
	"yang/browse"
)

//export
type BrowserHandle interface {
	browse.Browser
}

type c_browser struct {
	module *yang.Module
	enter_impl C.yangc2_browse_enter_impl
	iterate_impl C.yangc2_browse_iterate_impl
	read_impl C.yangc2_browse_read_impl
	edit_impl C.yangc2_browse_edit_impl
	choose_impl C.yangc2_browse_choose_impl
	exit_impl C.yangc2_browse_exit_impl
	root_selector_impl C.yangc2_browse_root_selector_impl
	browser_hnd unsafe.Pointer
}

//export
type ModuleHandle interface {
}

//export yangc2_load_module
func yangc2_load_module(
		enter_impl C.yangc2_browse_enter_impl,
		iterate_impl C.yangc2_browse_iterate_impl,
		read_impl C.yangc2_browse_read_impl,
		edit_impl C.yangc2_browse_edit_impl,
		choose_impl C.yangc2_browse_choose_impl,
		exit_impl C.yangc2_browse_exit_impl,
		root_selector_impl C.yangc2_browse_root_selector_impl,
		browser_hnd unsafe.Pointer,
        rs ResourceHandle,
        resourceId *C.char,
    ) ModuleHandle {

	var yang_module *yang.Module
	var module *yang.Module
	var err error

	module, err = yang.LoadModule(rs, C.GoString(resourceId))
	if (err != nil) {
		return nil
	}
	local_browser := browse.YangBrowser{Module:module, Meta:yang_module}

	remote_browser := yangc2_new_browser(
		enter_impl,
		iterate_impl,
		read_impl,
		edit_impl,
		choose_impl,
		exit_impl,
		module,
		root_selector_impl,
		browser_hnd,
	)

	to, err := remote_browser.RootSelector()
	if err == nil {
		from, err := local_browser.RootSelector()
		if err == nil {
			err = browse.Insert(from, to)
			if err != nil {
				// TODO: add module reference to driver so GC doesn't claim it
				return module
			}
		}
	}

	return nil
}

//export yangc2_new_browser
func yangc2_new_browser(
	enter_impl C.yangc2_browse_enter_impl,
	iterate_impl C.yangc2_browse_iterate_impl,
	read_impl C.yangc2_browse_read_impl,
	edit_impl C.yangc2_browse_edit_impl,
	choose_impl C.yangc2_browse_choose_impl,
	exit_impl C.yangc2_browse_exit_impl,
	module ModuleHandle,
	root_selector_impl C.yangc2_browse_root_selector_impl,
	browser_hnd unsafe.Pointer) (BrowserHandle) {

	return &c_browser{
		enter_impl:enter_impl,
		iterate_impl:iterate_impl,
		read_impl:read_impl,
		edit_impl:edit_impl,
		choose_impl:choose_impl,
		exit_impl:exit_impl,
		module: module.(*yang.Module),
		root_selector_impl: root_selector_impl,
		browser_hnd: browser_hnd,
	}
}

func (cb *c_browser) RootSelector() (s *browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	root_selection_hnd := C.yangc2_browse_root_selector(cb.root_selector_impl, cb.browser_hnd, errPtr)
	s, err = cb.c_selection(root_selection_hnd)
	s.Meta = cb.module
	return
}

func (cb *c_browser) c_selection(selection_hnd unsafe.Pointer) (*browse.Selection, error) {
	s := &browse.Selection{}
	s.Enter = func() (child *browse.Selection, err error) {
		return cb.c_enter(s, selection_hnd)
	}
	s.Iterate = func(keys []string, first bool) (bool, error) {
		return cb.c_iterate(s, selection_hnd, keys, first)
	}
	s.ReadValue = func(val *browse.Value) error {
		return cb.c_read(s, selection_hnd, val)
	}
	s.Edit = func(op browse.Operation, val *browse.Value) error {
		return cb.c_edit(s, selection_hnd, op, val)
	}
	s.Exit = func() error {
		return cb.c_exit(s, selection_hnd)
	}
	s.Choose = func(m *yang.Choice) (yang.Meta, error) {
		return cb.c_choose(s, selection_hnd, m)
	}
	return s, nil
}

func (cb *c_browser) c_iterate(s *browse.Selection, selection_hnd unsafe.Pointer, keys []string, first bool) (hasMore bool, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	var c_encoded_keys *C.char
	if len(keys) > 0 {
		c_encoded_keys = C.CString(strings.Join(keys, " "))
	}
	var c_first C.short
	if first {
		c_first = C.short(1)
	} else {
		c_first = C.short(0)

	}
	has_more := C.yangc2_browse_iterate(cb.enter_impl, selection_hnd, ident, c_encoded_keys, c_first, errPtr)

	return has_more > 0, err
}

func (cb *c_browser) c_enter(s *browse.Selection, selection_hnd unsafe.Pointer) (child *browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	c_found := C.short(0)
	ident := C.CString(s.Position.GetIdent())
	child_hnd := C.yangc2_browse_enter(cb.enter_impl, selection_hnd, ident, &c_found, errPtr)
	if c_found > 0 {
		s.Found = true
	}
	if child_hnd != nil && err == nil {
		child, err = cb.c_selection(child_hnd)
	}
	return
}

func (cb *c_browser) c_read(s *browse.Selection, selection_hnd unsafe.Pointer, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	var c_val C.struct_yangc2_browse_value
	C.yangc2_browse_read(cb.read_impl, selection_hnd, ident, &c_val, errPtr)
	switch c_val.val_type {
	case C.enum_yangc2_browse_value_type(C.STRING):
		val.Str = C.GoString(c_val.str)
	case C.enum_yangc2_browse_value_type(C.INT32):
		val.Int = int(c_val.int32)
	case C.enum_yangc2_browse_value_type(C.BOOLEAN):
		if c_val.boolean > C_FALSE {
			val.Bool = true
		} else {
			// nop
			val.Bool = false
		}
	}

	return
}

const (
	C_TRUE = C.short(1)
	C_FALSE = C.short(0)
)

func (cb *c_browser) c_edit(s *browse.Selection, selection_hnd unsafe.Pointer, op browse.Operation, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	var c_val C.struct_yangc2_browse_value
	dt := s.Position.(yang.HasDataType).GetDataType()
	switch dt.Ident {
	case "string":
		c_val.val_type = C.enum_yangc2_browse_value_type(C.STRING)
		c_val.str = C.CString(val.Str)
	case "int32":
		c_val.val_type = C.enum_yangc2_browse_value_type(C.INT32)
		c_val.int32 = C.int(val.Int)
	case "boolean":
		c_val.val_type = C.enum_yangc2_browse_value_type(C.BOOLEAN)
		if val.Bool {
			c_val.boolean = C_TRUE
		} else {
			c_val.boolean = C_FALSE
		}
	default:
		return &driverError{"Unsupported type"}
	}
	C.yangc2_browse_edit(cb.read_impl, selection_hnd, ident, C.int(op), &c_val, errPtr)
	return
}

func (cb *c_browser) c_choose(s *browse.Selection, selection_hnd unsafe.Pointer, choice *yang.Choice) (resolved yang.Meta, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	resolved_case := C.yangc2_browse_choose(cb.choose_impl, selection_hnd, ident, errPtr)
	if err == nil {
		ccase := choice.GetCase(C.GoString(resolved_case))
		resolved = ccase.GetFirstMeta()
	}
	return
}

func (cb *c_browser) c_exit(s *browse.Selection, selection_hnd unsafe.Pointer) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	C.yangc2_browse_exit(cb.choose_impl, selection_hnd, ident, errPtr)
	return
}
