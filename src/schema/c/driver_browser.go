package c

// #include "conf2/browse.h"
// extern void *conf2_browse_root_selector(conf2_browse_root_selector_impl impl_func, void *browser_handle, void *browse_err);
// extern void *conf2_browse_enter(conf2_browse_enter_impl impl_func, void *selection_handle, char *ident, short *found, void *browse_err);
// extern short conf2_browse_iterate(conf2_browse_iterate_impl impl_func, void *selection_handle, void *key_data, int key_data_len, short first, void *browse_err);
// extern void *conf2_browse_read(conf2_browse_read_impl impl_func, void *selection_handle, char *ident, void **val_data_ptr, int* val_data_len_ptr, void *browse_err);
// extern void conf2_browse_edit(conf2_browse_edit_impl impl_func, void *selection_handle, char *ident, int op, void *val_data, int val_data_len, void *browse_err);
// extern char *conf2_browse_choose(conf2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern void conf2_browse_exit(conf2_browse_exit_impl impl_func, void *selection_handle, char *ident, void *browse_err);
import "C"

import (
	"schema"
	"unsafe"
	"schema/browse"
	"schema/yang"
	"schema/comm"
	"fmt"
)

type apiBrowser struct {
	module *schema.Module
	enter_impl C.conf2_browse_enter_impl
	iterate_impl C.conf2_browse_iterate_impl
	read_impl C.conf2_browse_read_impl
	edit_impl C.conf2_browse_edit_impl
	choose_impl C.conf2_browse_choose_impl
	exit_impl C.conf2_browse_exit_impl
	root_selector_impl C.conf2_browse_root_selector_impl
	browser_hnd *ApiHandle
}

//export conf2_load_module
func conf2_load_module(
		enter_impl C.conf2_browse_enter_impl,
		iterate_impl C.conf2_browse_iterate_impl,
		read_impl C.conf2_browse_read_impl,
		edit_impl C.conf2_browse_edit_impl,
		choose_impl C.conf2_browse_choose_impl,
		exit_impl C.conf2_browse_exit_impl,
		root_selector_impl C.conf2_browse_root_selector_impl,
		browser_hnd_id unsafe.Pointer,
        stream_source_hnd_id unsafe.Pointer,
        resourceId *C.char,
    ) unsafe.Pointer {

	stream_source_hnd, stream_source_found := GoHandles()[stream_source_hnd_id]
	if !stream_source_found {
		panic(fmt.Sprint("Stream source not found", stream_source_hnd_id))
	}
	stream_source := stream_source_hnd.Data.(schema.StreamSource)
	defer schema.CloseResource(stream_source)

	var module *schema.Module
	var err error

	module, err = yang.LoadModule(stream_source, C.GoString(resourceId))
	if (err != nil) {
		fmt.Println("Error loading module", err.Error())
		return nil
	}
	module_browser := browse.NewSchemaBrowser(module, false)

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

	var to browse.Selection
	to, _, err = apiModuleBrowser.RootSelector()
	defer schema.CloseResource(to)
	if err == nil {
		var from browse.Selection
		var state *browse.WalkState
		from, state, err = module_browser.RootSelector()
		defer schema.CloseResource(from)
		if err == nil {
			err = browse.Insert(state, from, to, browse.WalkAll())
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

//export conf2_new_browser
func conf2_new_browser(
	enter_impl C.conf2_browse_enter_impl,
	iterate_impl C.conf2_browse_iterate_impl,
	read_impl C.conf2_browse_read_impl,
	edit_impl C.conf2_browse_edit_impl,
	choose_impl C.conf2_browse_choose_impl,
	exit_impl C.conf2_browse_exit_impl,
    root_selector_impl C.conf2_browse_root_selector_impl,
	module_hnd_id unsafe.Pointer,
	browser_hnd_id unsafe.Pointer) unsafe.Pointer {

	module_hnd, module_found := GoHandles()[module_hnd_id]
	if ! module_found {
		panic(fmt.Sprintf("Module not found %p\n", module_hnd_id))
		return nil
	}
	module := module_hnd.Data.(*schema.Module)
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

func (cb *apiBrowser) RootSelector() (s browse.Selection, state *browse.WalkState, err error) {
	errPtr := unsafe.Pointer(&err)
	selector := &apiSelector{browser:cb}
	root_selection_hnd_id := C.conf2_browse_root_selector(cb.root_selector_impl, cb.browser_hnd.ID, errPtr)
	root_selection_hnd, found := ApiHandles()[root_selection_hnd_id]
	if ! found {
		panic(fmt.Sprint("Root selector handle not found", root_selection_hnd_id))
	}
	s, err = selector.selection(root_selection_hnd)
	state = browse.NewWalkState(cb.module)
	return
}

func (cb *apiBrowser) Module() *schema.Module {
	return cb.module
}

func (cb *apiBrowser) Close() error {
	return cb.browser_hnd.Close()
}

type apiSelector struct {
	browser *apiBrowser
}

func (cb *apiSelector) selection(selectionHnd *ApiHandle) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		return cb.enter(meta, selectionHnd)
	}
	s.OnNext = func(state *browse.WalkState, meta *schema.List, keys []*browse.Value, first bool) (bool, error) {
		return cb.iterate(s, selectionHnd, keys, first)
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		return cb.read(meta, selectionHnd)
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, val *browse.Value) error {
		return cb.edit(meta, selectionHnd, op, val)
	}
	s.OnUnselect = func(state *browse.WalkState, meta schema.MetaList) error {
		return cb.exit(meta, selectionHnd)
	}
	s.OnChoose = func(state *browse.WalkState, m *schema.Choice) (schema.Meta, error) {
		return cb.choose(state.Position(), selectionHnd, m)
	}
	s.Resource = selectionHnd
	return s, nil
}


func (cb *apiSelector) iterate(s browse.Selection, selectionHnd *ApiHandle, key []*browse.Value, first bool) (hasMore bool, err error) {
	errPtr := unsafe.Pointer(&err)
	var c_encoded_key unsafe.Pointer
	var encoded_key_len C.int
	if len(key) > 0 {
		w := comm.NewWriter()
		w.WriteValues(key)
		data := w.Data()
		encoded_key_len = C.int(len(data))
		c_encoded_key = NewGoHandleByteArray(data).ID
	}
	var c_first C.short
	if first {
		c_first = C.short(1)
	} else {
		c_first = C.short(0)

	}
	has_more := C.conf2_browse_iterate(cb.browser.iterate_impl, selectionHnd.ID, c_encoded_key, encoded_key_len, c_first, errPtr)

	return has_more > 0, err
}

func (cb *apiSelector) enter(meta schema.MetaList, selectionHnd *ApiHandle) (child browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	c_found := C.short(0)
	ident := encodeIdent(meta)
	child_hnd_id := C.conf2_browse_enter(cb.browser.enter_impl, selectionHnd.ID, ident, &c_found, errPtr)
	if c_found == 0 {
		return nil, nil
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

func encodeIdent(position schema.Meta) *C.char {
	if ccase, isCase := position.GetParent().(*schema.ChoiceCase); isCase {
		path := fmt.Sprintf("%s/%s/%s", ccase.GetParent().GetIdent(), ccase.GetIdent(), position.GetIdent())
		return C.CString(path)
	}
	return C.CString(position.GetIdent())
}

func (cb *apiSelector) read(meta schema.HasDataType, selectionHnd *ApiHandle) (val *browse.Value, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := encodeIdent(meta)
	var data_ptr unsafe.Pointer
	var data_len C.int
	value_hnd_id := C.conf2_browse_read(cb.browser.read_impl, selectionHnd.ID, ident, &data_ptr, &data_len, errPtr)
	if data_ptr != nil && data_len > 0 {
		data := C.GoBytes(data_ptr, data_len)
		r := comm.NewReader(data)
		expectedType := meta.GetDataType()
		val, err = r.ReadValue(expectedType)
		if value_hnd_id != nil {
			if value_hnd, found := ApiHandles()[value_hnd_id]; found {
				value_hnd.Close()
			}
		}
	}

	return val, err
}

func (cb *apiSelector) edit(meta schema.Meta, selectionHnd *ApiHandle, op browse.Operation, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	var handle *GoHandle
	ident := encodeIdent(meta)
	var data_ptr unsafe.Pointer
	var data_len C.int
	if val != nil {
		w := comm.NewWriter()
		w.WriteValue(val)
		data := w.Data()
		data_len = C.int(len(data))
		handle = NewGoHandleByteArray(data)
		data_ptr = handle.ID
	}

	C.conf2_browse_edit(cb.browser.edit_impl, selectionHnd.ID, ident, C.int(op), data_ptr, data_len, errPtr)

	if handle != nil {
		handle.Close()
	}

	return
}

func (cb *apiSelector) choose(meta schema.Meta, selectionHnd *ApiHandle, choice *schema.Choice) (resolved schema.Meta, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(meta.GetIdent())
	resolved_case := C.conf2_browse_choose(cb.browser.choose_impl, selectionHnd.ID, ident, errPtr)
	if err == nil {
		ccase := choice.GetCase(C.GoString(resolved_case))
		resolved = ccase.GetFirstMeta()
	}
	return
}

func (cb *apiSelector) exit(meta schema.MetaList, selectionHnd *ApiHandle) (err error) {
	errPtr := unsafe.Pointer(&err)
	ident := encodeIdent(meta)
	C.conf2_browse_exit(cb.browser.exit_impl, selectionHnd.ID, ident, errPtr)
	return
}
