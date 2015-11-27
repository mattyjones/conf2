package c

// #include "conf2/browse.h"
// extern void *conf2_browse_selector(conf2_browse_selector_impl impl_func, void *browser_handle, char *path, void *browse_err);
// extern void *conf2_browse_enter(conf2_browse_enter_impl impl_func, void *selection_handle, char *ident, short create, void *browse_err);
// extern void *conf2_browse_iterate(conf2_browse_iterate_impl impl_func, void *selection_handle, short create, void *key_data, int key_data_len, short first, void *browse_err);
// extern void *conf2_browse_read(conf2_browse_read_impl impl_func, void *selection_handle, char *ident, void **val_data_ptr, int* val_data_len_ptr, void *browse_err);
// extern void conf2_browse_edit(conf2_browse_edit_impl impl_func, void *selection_handle, char *ident, void *val_data, int val_data_len, void *browse_err);
// extern char *conf2_browse_choose(conf2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern void conf2_browse_event(conf2_browse_event_impl impl_func, void *selection_handle, int event_id, void *browse_err);
import "C"

import (
	"schema"
	"unsafe"
	"data"
	"schema/yang"
	"comm"
	"fmt"
)

type apiBrowser struct {
	module *schema.Module
	enter_impl C.conf2_browse_enter_impl
	iterate_impl C.conf2_browse_iterate_impl
	read_impl C.conf2_browse_read_impl
	edit_impl C.conf2_browse_edit_impl
	choose_impl C.conf2_browse_choose_impl
	event_impl C.conf2_browse_event_impl
	selector_impl C.conf2_browse_selector_impl
	browser_hnd *ApiHandle
}

//export conf2_load_module
func conf2_load_module(
		enter_impl C.conf2_browse_enter_impl,
		iterate_impl C.conf2_browse_iterate_impl,
		read_impl C.conf2_browse_read_impl,
		edit_impl C.conf2_browse_edit_impl,
		choose_impl C.conf2_browse_choose_impl,
		event_impl C.conf2_browse_event_impl,
		selector_impl C.conf2_browse_selector_impl,
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
	module_browser := data.NewSchemaData(module, false)

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
		event_impl:event_impl,
		module: module,
		selector_impl: selector_impl,
		browser_hnd: browser_hnd,
	}

	var to, from *data.Selection
	rootPath := data.NewPath("")
	to, err = apiModuleBrowser.Selector(rootPath)
	defer schema.CloseResource(to)
	if err == nil {
		from, err = module_browser.Selector(rootPath)
		defer schema.CloseResource(from)
		if err == nil {
			err = data.Insert(from, to)
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
	event_impl C.conf2_browse_event_impl,
    selector_impl C.conf2_browse_selector_impl,
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
		event_impl:event_impl,
		module: module,
		selector_impl: selector_impl,
		browser_hnd: browser_hnd,
	}
	return NewGoHandle(browser).ID
}

func (cb *apiBrowser) Selector(path *data.Path) (sel *data.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	selector := &apiSelector{browser:cb}
	cstr_path := C.CString(path.URL)
	root_selection_hnd_id := C.conf2_browse_selector(cb.selector_impl, cb.browser_hnd.ID, cstr_path, errPtr)
	root_selection_hnd, found := ApiHandles()[root_selection_hnd_id]
	if ! found {
		panic(fmt.Sprint("Root selector handle not found", root_selection_hnd_id))
	}
	var n data.Node
	n, err = selector.selection(root_selection_hnd)

	// TODO: Need to pass path and allow driver to walk and return state, otherwise destination
	// cannot fast-forward based on path
	return data.WalkPath(data.NewSelection(n, cb.module), path)
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

func (cb *apiSelector) selection(selectionHnd *ApiHandle) (data.Node, error) {
	s := &data.MyNode{}
	s.OnSelect = func(state *data.Selection, meta schema.MetaList, new bool) (child data.Node, err error) {
		return cb.enter(meta, selectionHnd, new)
	}
	s.OnNext = func(state *data.Selection, meta *schema.List, new bool, keys []*data.Value, first bool) (data.Node, error) {
		return cb.iterate(s, selectionHnd, new, keys, first)
	}
	s.OnRead = func(state *data.Selection, meta schema.HasDataType) (*data.Value, error) {
		return cb.read(meta, selectionHnd)
	}
	s.OnWrite = func(state *data.Selection, meta schema.HasDataType, val *data.Value) error {
		return cb.edit(meta, selectionHnd, val)
	}
	s.OnEvent = func(state *data.Selection, e data.Event) error {
		return cb.event(selectionHnd, e)
	}
	s.OnChoose = func(state *data.Selection, m *schema.Choice) (schema.Meta, error) {
		return cb.choose(state.State.Position(), selectionHnd, m)
	}
	s.Resource = selectionHnd
	return s, nil
}

func (cb *apiSelector) iterate(s data.Node, selectionHnd *ApiHandle, new bool, key []*data.Value, first bool) (data.Node, error) {
	var err error
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
	var c_create C.short
	if new {
		c_create = C.short(1)
	} else {
		c_create = C.short(0)
	}

	child_hnd_id := C.conf2_browse_iterate(cb.browser.iterate_impl, selectionHnd.ID, c_create, c_encoded_key, encoded_key_len, c_first, errPtr)
	if child_hnd_id == nil || err != nil {
		return nil, err
	}
	child_hnd, found := ApiHandles()[child_hnd_id]
	if ! found {
		panic(fmt.Sprint("Enter selector handle not found", child_hnd_id))
	}
	return cb.selection(child_hnd)
}

func (cb *apiSelector) enter(meta schema.MetaList, selectionHnd *ApiHandle, new bool) (child data.Node, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := encodeIdent(meta)
	var c_create C.short
	if new {
		c_create = C.short(1)
	} else {
		c_create = C.short(0)
	}
	child_hnd_id := C.conf2_browse_enter(cb.browser.enter_impl, selectionHnd.ID, ident, c_create, errPtr)
	if child_hnd_id == nil || err != nil {
		return nil, err
	}
	child_hnd, found := ApiHandles()[child_hnd_id]
	if ! found {
		panic(fmt.Sprint("Enter selector handle not found", child_hnd_id))
	}
	return cb.selection(child_hnd)
}

func encodeIdent(position schema.Meta) *C.char {
	if ccase, isCase := position.GetParent().(*schema.ChoiceCase); isCase {
		path := fmt.Sprintf("%s/%s/%s", ccase.GetParent().GetIdent(), ccase.GetIdent(), position.GetIdent())
		return C.CString(path)
	}
	return C.CString(position.GetIdent())
}

func (cb *apiSelector) read(meta schema.HasDataType, selectionHnd *ApiHandle) (val *data.Value, err error) {
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

func (cb *apiSelector) edit(meta schema.HasDataType, selectionHnd *ApiHandle, val *data.Value) (err error) {
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

	C.conf2_browse_edit(cb.browser.edit_impl, selectionHnd.ID, ident, data_ptr, data_len, errPtr)

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

func (cb *apiSelector) event(selectionHnd *ApiHandle, e data.Event) (err error) {
	errPtr := unsafe.Pointer(&err)
	C.conf2_browse_event(cb.browser.event_impl, selectionHnd.ID, C.int(e), errPtr)
	return
}
