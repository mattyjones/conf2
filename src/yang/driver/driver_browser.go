package driver

// #include "yang-c2/browse.h"
// extern void* yangc2_browse_root_selector(yangc2_browse_root_selector_impl impl_func, void *browser_handle, void *browse_err);
// extern void* yangc2_browse_enter(yangc2_browse_enter_impl impl_func, void *selection_handle, char *ident, short *found, void *browse_err);
// extern short yangc2_browse_iterate(yangc2_browse_iterate_impl impl_func, void *selection_handle, char *encodedKeys, short first, void *browse_err);
// extern void yangc2_browse_read(yangc2_browse_read_impl impl_func, void *selection_handle, char *ident, struct yangc2_browse_value *val, void *browse_err);
// extern void yangc2_browse_edit(yangc2_browse_edit_impl impl_func, void *selection_handle, char *ident, int op, struct yangc2_browse_value *val, void *browse_err);
// extern char *yangc2_browse_choose(yangc2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern void yangc2_browse_exit(yangc2_browse_exit_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern char** yangc2_cstrslice_as_strlist(void *cstr_slice);
// extern int* yangc2_cintslice_as_intlist(void *cint_slice);
// extern short* yangc2_cboolslice_as_boollist(void *cbool_slice);
import "C"

import (
	"yang"
	"unsafe"
	"strings"
	"yang/browse"
	"fmt"
	"bytes"
	"bufio"
)

type apiBrowser struct {
	module *yang.Module
	enter_impl C.yangc2_browse_enter_impl
	iterate_impl C.yangc2_browse_iterate_impl
	read_impl C.yangc2_browse_read_impl
	edit_impl C.yangc2_browse_edit_impl
	choose_impl C.yangc2_browse_choose_impl
	exit_impl C.yangc2_browse_exit_impl
	root_selector_impl C.yangc2_browse_root_selector_impl
	browser_hnd *ApiHandle
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
		browser_hnd_id unsafe.Pointer,
        stream_source_hnd_id unsafe.Pointer,
        resourceId *C.char,
    ) unsafe.Pointer {

	stream_source_hnd, stream_source_found := GoHandles()[stream_source_hnd_id]
	if !stream_source_found {
		panic(fmt.Sprint("Stream source not found", stream_source_hnd_id))
	}
	stream_source := stream_source_hnd.Data.(yang.StreamSource)
	defer stream_source.Close()

	var module *yang.Module
	var err error

	module, err = yang.LoadModule(stream_source, C.GoString(resourceId))
	if (err != nil) {
		fmt.Println("Error loading module", err.Error())
		return nil
	}
	module_browser := browse.NewYangBrowser(module)

	browser_hnd, browser_found := ApiHandles()[browser_hnd_id]
	if ! browser_found {
		panic(fmt.Sprint("Browser not found", browser_hnd))
	}
	apiModuleBrowser := &apiBrowser{
		enter_impl:enter_impl,
		iterate_impl:iterate_impl,
		read_impl:read_impl,
		edit_impl:edit_impl,
		choose_impl:choose_impl,
		exit_impl:exit_impl,
		module: module,
		root_selector_impl: root_selector_impl,
		browser_hnd: browser_hnd,
	}

	var to *browse.Selection
	to, err = apiModuleBrowser.RootSelector()
	defer to.Close()
	if err == nil {
		var from *browse.Selection
		from, err = module_browser.RootSelector()
		defer from.Close()
		if err == nil {
			err = browse.Insert(from, to)
			if err == nil {
				moduleHnd := NewGoHandle(module)
				return moduleHnd.ID
			}
		}
	}
	if (err != nil) {
		fmt.Println("Error sending module", err.Error())
		return nil
	}

	// TODO: Find a way to return an error and not just nil
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
    root_selector_impl C.yangc2_browse_root_selector_impl,
	module_hnd_id unsafe.Pointer,
	browser_hnd_id unsafe.Pointer) unsafe.Pointer {

	module_hnd, module_found := GoHandles()[module_hnd_id]
	if ! module_found {
		panic(fmt.Sprintf("Module not found %p\n", module_hnd_id))
		return nil
	}
	module := module_hnd.Data.(*yang.Module)
	browser_hnd, found_browser := ApiHandles()[browser_hnd_id]
	if ! found_browser {

		return nil
	}

	browser := &apiBrowser{
		enter_impl:enter_impl,
		iterate_impl:iterate_impl,
		read_impl:read_impl,
		edit_impl:edit_impl,
		choose_impl:choose_impl,
		exit_impl:exit_impl,
		module: module,
		root_selector_impl: root_selector_impl,
		browser_hnd: browser_hnd,
	}
	return NewGoHandle(browser).ID
}

func (cb *apiBrowser) RootSelector() (s *browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	selector := &apiSelector{browser:cb}
	root_selection_hnd_id := C.yangc2_browse_root_selector(cb.root_selector_impl, cb.browser_hnd.ID, errPtr)
	root_selection_hnd, found := ApiHandles()[root_selection_hnd_id]
	if ! found {
		panic(fmt.Sprint("Root selector handle not found", root_selection_hnd_id))
	}
	s, err = selector.selection(root_selection_hnd)
	s.Meta = cb.module
	s.Resource = root_selection_hnd
	return
}

func (cb *apiBrowser) Module() *yang.Module {
	return cb.module
}

func (cb *apiBrowser) Close() error {
	return cb.browser_hnd.Close()
}

type apiSelector struct {
	browser *apiBrowser
}

func (cb *apiSelector) selection(selectionHnd *ApiHandle) (*browse.Selection, error) {
	s := &browse.Selection{}
	s.Enter = func() (child *browse.Selection, err error) {
		return cb.enter(s, selectionHnd)
	}
	s.Iterate = func(keys []string, first bool) (bool, error) {
		return cb.iterate(s, selectionHnd, keys, first)
	}
	s.ReadValue = func(val *browse.Value) error {
		return cb.read(s, selectionHnd, val)
	}
	s.Edit = func(op browse.Operation, val *browse.Value) error {
		return cb.edit(s, selectionHnd, op, val)
	}
	s.Exit = func() error {
		return cb.exit(s, selectionHnd)
	}
	s.Choose = func(m *yang.Choice) (yang.Meta, error) {
		return cb.choose(s, selectionHnd, m)
	}
	return s, nil
}

func (cb *apiSelector) iterate(s *browse.Selection, selectionHnd *ApiHandle, keys []string, first bool) (hasMore bool, err error) {
	errPtr := unsafe.Pointer(&err)
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
	has_more := C.yangc2_browse_iterate(cb.browser.iterate_impl, selectionHnd.ID, c_encoded_keys, c_first, errPtr)

	return has_more > 0, err
}

func (cb *apiSelector) enter(s *browse.Selection, selectionHnd *ApiHandle) (child *browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	c_found := C.short(0)
	ident := encodeIdent(s.Position)
	child_hnd_id := C.yangc2_browse_enter(cb.browser.enter_impl, selectionHnd.ID, ident, &c_found, errPtr)
	if c_found > 0 {
		s.Found = true
	}
	if child_hnd_id != nil && err == nil {
		child_hnd, found := ApiHandles()[child_hnd_id]
		if ! found {
			panic(fmt.Sprint("Enter selector handle not found", child_hnd_id))
		}
		child, err = cb.selection(child_hnd)
	}
	return
}

func encodeIdent(position yang.Meta) *C.char {
	if ccase, isCase := position.GetParent().(*yang.ChoiceCase); isCase {
		path := fmt.Sprintf("%s/%s/%s", ccase.GetParent().GetIdent(), ccase.GetIdent(), position.GetIdent())
		return C.CString(path)
	}
	return C.CString(position.GetIdent())
}

func (cb *apiSelector) read_cstrs(cstr_list []*C.char, len int) []string {

}

func (cb *apiSelector) read(s *browse.Selection, selectionHnd *ApiHandle, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := encodeIdent(s.Position)
	var c_val C.struct_yangc2_browse_value
	C.yangc2_browse_read(cb.browser.read_impl, selectionHnd.ID, ident, &c_val, errPtr)
	if c_val.is_list > 0 {
		val.IsList = true
		switch c_val.val_type {
		case C.enum_yangc2_browse_value_type(C.ENUMERATION):
			val.Strlist == read_cstrs(c_val.cstr_list, c_val.list_len)
			val.Intlist == read_ints(c_val.int_list, c_val.list_len)
			val.Type = yang.TYPE_ENUMERATION
		case C.enum_yangc2_browse_value_type(C.STRING):
			val.Strlist == read_cstrs(c_val.cstr_list, c_val.list_len)
			val.Type = yang.TYPE_STRING
		case C.enum_yangc2_browse_value_type(C.INT32):
			val.Intlist == read_ints(c_val.int_list, c_val.list_len)
			val.Type = yang.TYPE_INT32
		case C.enum_yangc2_browse_value_type(C.BOOLEAN):
			if c_val.boolean > C_FALSE {
				val.Bool = true
			} else {
				// nop
				val.Bool = false
			}
			val.Type = yang.TYPE_BOOLEAN

		}
	} else {
		switch c_val.val_type {
		case C.enum_yangc2_browse_value_type(C.ENUMERATION):
			val.Str = C.GoString(c_val.cstr)
			val.Int = int(c_val.int32)
		case C.enum_yangc2_browse_value_type(C.STRING):
			val.Str = C.GoString(c_val.cstr)
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
	}

	return
}

const (
	C_TRUE = C.short(1)
	C_FALSE = C.short(0)
)

func leafListValue(val *browse.Value) (*C.struct_yangc2_browse_value, error) {
	var buff bytes.Buffer
	w := bufio.NewWriter(buff)
	var c_val C.struct_yangc2_browse_value
	c_val.is_list = C_TRUE
	switch val.Type.Ident {
	case "string":
		c_val.list_len = C.int(len(val.Strlist))
		strlist :=  make([]*C.char, len(val.Strlist))
		for i, s := range val.Strlist {
			strlist[i] = C.CString(s)
		}
		c_val.handle = NewGoHandle(strlist).ID
		c_val.val_type = C.enum_yangc2_browse_value_type(C.STRING)
		c_val.cstr_list = C.yangc2_cstrslice_as_strlist(unsafe.Pointer(&strlist))
	case "int32":
		c_val.list_len = C.int(len(val.Intlist))
		intlist := make([]C.int, len(val.Intlist))
		c_val.handle = NewGoHandle(intlist).ID
		for i, ival := range val.Intlist {
			intlist[i] = C.int(ival)
		}
		c_val.val_type = C.enum_yangc2_browse_value_type(C.INT32)
		c_val.int_list = C.yangc2_cintslice_as_intlist(unsafe.Pointer(&intlist))
	case "boolean":
		// TODO: could make this smaller using bit field
		c_val.list_len = C.int(len(val.Boollist))
		boollist := make([]C.short, len(val.Boollist))
		c_val.handle = NewGoHandle(boollist).ID
		for i, bval := range val.Boollist {
			if bval {
				boollist[i] = C_TRUE
			} else {
				boollist[i] = C_FALSE
			}
		}
		c_val.val_type = C.enum_yangc2_browse_value_type(C.BOOLEAN)
		c_val.bool_list = C.yangc2_cboolslice_as_boollist(unsafe.Pointer(&boollist))
	default:
		return nil, &driverError{"Unsupported type"}
	}
	return &c_val, nil
}

func leafValue(val *browse.Value) (*C.struct_yangc2_browse_value, error) {
	var c_val C.struct_yangc2_browse_value
	switch val.Type.Ident {
	case "string":
		c_val.val_type = C.enum_yangc2_browse_value_type(C.STRING)
		c_val.cstr = C.CString(val.Str)
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
		return nil, &driverError{"Unsupported type"}
	}
	return &c_val, nil
}

func (cb *apiSelector) edit(s *browse.Selection, selectionHnd *ApiHandle, op browse.Operation, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	var ident *C.char;
	var c_val *C.struct_yangc2_browse_value
	if s.Position != nil {
		ident = encodeIdent(s.Position)
	}
	if val != nil {
		if val.IsList {
			c_val, err = leafListValue(val)
		} else {
			c_val, err = leafValue(val)
		}
	}
	C.yangc2_browse_edit(cb.browser.edit_impl, selectionHnd.ID, ident, C.int(op), c_val, errPtr)

	if val != nil && c_val.handle != nil {
		GoHandles()[c_val.handle].Close()
	}
	return
}

func (cb *apiSelector) choose(s *browse.Selection, selectionHnd *ApiHandle, choice *yang.Choice) (resolved yang.Meta, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.Position.GetIdent())
	resolved_case := C.yangc2_browse_choose(cb.browser.choose_impl, selectionHnd.ID, ident, errPtr)
	if err == nil {
		ccase := choice.GetCase(C.GoString(resolved_case))
		resolved = ccase.GetFirstMeta()
	}
	return
}

func (cb *apiSelector) exit(s *browse.Selection, selectionHnd *ApiHandle) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := encodeIdent(s.Position)
	C.yangc2_browse_exit(cb.browser.exit_impl, selectionHnd.ID, ident, errPtr)
	return
}
