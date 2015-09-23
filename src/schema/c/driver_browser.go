package c

// #include "conf2/browse.h"
// extern void* conf2_browse_root_selector(conf2_browse_root_selector_impl impl_func, void *browser_handle, void *browse_err);
// extern void* conf2_browse_enter(conf2_browse_enter_impl impl_func, void *selection_handle, char *ident, short *found, void *browse_err);
// extern short conf2_browse_iterate(conf2_browse_iterate_impl impl_func, void *selection_handle, char *encodedKeys, short first, void *browse_err);
// extern void conf2_browse_read(conf2_browse_read_impl impl_func, void *selection_handle, char *ident, struct conf2_value *val, void *browse_err);
// extern void conf2_browse_edit(conf2_browse_edit_impl impl_func, void *selection_handle, char *ident, int op, struct conf2_value *val, void *browse_err);
// extern char *conf2_browse_choose(conf2_browse_choose_impl impl_func, void *selection_handle, char *ident, void *browse_err);
// extern void conf2_browse_exit(conf2_browse_exit_impl impl_func, void *selection_handle, char *ident, void *browse_err);
import "C"

import (
	"schema"
	"unsafe"
	"strings"
	"schema/browse"
	"schema/yang"
	"fmt"
	"bytes"
	"bufio"
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
	to, err = apiModuleBrowser.RootSelector()
	defer schema.CloseResource(to)
	if err == nil {
		var from browse.Selection
		from, err = module_browser.RootSelector()
		defer schema.CloseResource(from)
		if err == nil {
			err = browse.Insert(from, to, browse.NewExhaustiveController())
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

func (cb *apiBrowser) RootSelector() (s browse.Selection, err error) {
	errPtr := unsafe.Pointer(&err)
	selector := &apiSelector{browser:cb}
	root_selection_hnd_id := C.conf2_browse_root_selector(cb.root_selector_impl, cb.browser_hnd.ID, errPtr)
	root_selection_hnd, found := ApiHandles()[root_selection_hnd_id]
	if ! found {
		panic(fmt.Sprint("Root selector handle not found", root_selection_hnd_id))
	}
	s, err = selector.selection(root_selection_hnd)
	s.WalkState().Meta = cb.module
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
	s.OnSelect = func(meta schema.MetaList) (child browse.Selection, err error) {
		return cb.enter(meta, selectionHnd)
	}
	s.OnNext = func(keys []*browse.Value, first bool) (bool, error) {
		return cb.iterate(s, selectionHnd, keys, first)
	}
	s.OnRead = func(meta schema.HasDataType) (*browse.Value, error) {
		return cb.read(meta, selectionHnd)
	}
	s.OnWrite = func(meta schema.Meta, op browse.Operation, val *browse.Value) error {
		return cb.edit(s, selectionHnd, op, val)
	}
	s.OnUnselect = func(meta schema.MetaList) error {
		return cb.exit(meta, selectionHnd)
	}
	s.OnChoose = func(m *schema.Choice) (schema.Meta, error) {
		return cb.choose(s, selectionHnd, m)
	}
	s.Resource = selectionHnd
	return s, nil
}

func encodeKey(key []*browse.Value) []byte {

}

func (cb *apiSelector) iterate(s browse.Selection, selectionHnd *ApiHandle, key []*browse.Value, first bool) (hasMore bool, err error) {
	errPtr := unsafe.Pointer(&err)
	var c_encoded_keys *C.char
	if len(key) > 0 {
		c_encoded_key = encodeKey(key) C.CString(strings.Join(keys, " "))
	}
	var c_first C.short
	if first {
		c_first = C.short(1)
	} else {
		c_first = C.short(0)

	}
	has_more := C.conf2_browse_iterate(cb.browser.iterate_impl, selectionHnd.ID, c_encoded_key, c_first, errPtr)

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
	var c_val C.struct_conf2_value
	C.conf2_browse_read(cb.browser.read_impl, selectionHnd.ID, ident, &c_val, errPtr)
	valType := meta.GetDataType()
	c_val_format := int(c_val.format)
	if int(valType.Format) != c_val_format {
		e := fmt.Sprint("Invalid value, expected %d but got %d", valType.Format, c_val_format)
		return nil, &driverError{e}
	}
	val = &browse.Value{IsList: (c_val.is_list > 0) }
	if val.IsList {
		var data []byte
		data = C.GoBytes(c_val.data, c_val.data_len)
		listLen := int(c_val.list_len)

		switch val.Type.Format {
		case schema.FMT_ENUMERATION:
			intlist := getIntSlice(data, listLen)
			val.SetEnumList(intlist)
		case schema.FMT_STRING:
			val.Strlist = getStringSlice(data, listLen, int(c_val.data_len))
		case schema.FMT_INT32:
			val.Intlist = getIntSlice(data, listLen)
		case schema.FMT_BOOLEAN:
			val.Boollist = getBoolSlice(data, listLen)
		}
	} else {
		switch val.Type.Format {
		case schema.FMT_ENUMERATION:
			val.SetEnum(int(c_val.int32))
		case schema.FMT_STRING:
			val.Str = C.GoString(c_val.cstr)
		case schema.FMT_INT32:
			val.Int = int(c_val.int32)
		case schema.FMT_BOOLEAN:
			if c_val.boolean > C_FALSE {
				val.Bool = true
			}
		}
	}

	return val, err
}

const (
	C_TRUE = C.short(1)
	C_FALSE = C.short(0)
)

func leafListValue(val *browse.Value) (*GoHandle, *C.struct_conf2_value, error) {
	var c_val C.struct_conf2_value
	var w := NewCWriter()
	c_val.is_list, c_val.list_len = w.putValue(val)
	data := w.Bytes()

//	c_val.data, c_val.data_len =
//	c_val.format = uint32(val.Type.Format)
//	c_val.is_list = C_TRUE
//	switch val.Type.Format {
//	case schema.FMT_STRING:
//		c_val.list_len = C.int(len(val.Strlist))
//		for _, s := range val.Strlist {
//			putString(buffer, s)
//			buffer.WriteString(s)
//			buffer.WriteByte(CSTR_TERM)
//		}
//	case schema.FMT_INT32, schema.FMT_ENUMERATION:
//		bInt := [4]byte{}
//		c_val.list_len = C.int(len(val.Intlist))
//		for _, i := range val.Intlist {
//			putInt(&bInt, i)
//			buffer.Write(bInt[:])
//		}
//	case schema.FMT_BOOLEAN:
//		// TODO: Performance - could make this smaller using bit field
//		bShort := [2]byte{}
//		c_val.list_len = C.int(len(val.Boollist))
//		for _, b := range val.Boollist {
//			if b {
//				putShort(&bShort, C_TRUE)
//			} else {
//				putShort(&bShort, C_FALSE)
//			}
//			buffer.Write(bShort[:])
//		}
//	default:
//		return nil, nil, &driverError{"Unsupported type"}
//	}
//	buffer.Flush()
//	bytes := byteBuff.Bytes()
	if len(data) > 0 {
		c_val.data = unsafe.Pointer(&data[0])
		c_val.data_len = C.int(len(data))
	}
	c_val.format = uint32(val.Type.Format)
	// we return a handle to the bytes because nothing references the object leaving
	// this function and in theory could be GC'ed.  Handle can be released when
	// c_val is GC'ed
	return NewGoHandle(bytes), &c_val, nil
}

func leafValue(val *browse.Value) (*C.struct_conf2_value, error) {
	var c_val C.struct_conf2_value
	switch val.Type.Format {
	case schema.FMT_STRING:
		c_val.cstr = C.CString(val.Str)
	case schema.FMT_INT32, schema.FMT_ENUMERATION:
		c_val.int32 = C.int(val.Int)
	case schema.FMT_BOOLEAN:
		if val.Bool {
			c_val.boolean = C_TRUE
		} else {
			c_val.boolean = C_FALSE
		}
	default:
		return nil, &driverError{"Unsupported type"}
	}
	c_val.format = uint32(val.Type.Format)
	return &c_val, nil
}

func (cb *apiSelector) edit(s browse.Selection, selectionHnd *ApiHandle, op browse.Operation, val *browse.Value) (err error) {
	errPtr := unsafe.Pointer(&err)
	var ident *C.char;
	var c_val *C.struct_conf2_value
	var handle *GoHandle
	if s.WalkState().Position != nil {
		ident = encodeIdent(s.WalkState().Position)
	}
	if val != nil {
		if val.IsList {
			handle, c_val, err = leafListValue(val)
		} else {
			c_val, err = leafValue(val)
		}
	}

	C.conf2_browse_edit(cb.browser.edit_impl, selectionHnd.ID, ident, C.int(op), c_val, errPtr)

	if handle != nil {
		handle.Close()
	}

	return
}

func (cb *apiSelector) choose(s browse.Selection, selectionHnd *ApiHandle, choice *schema.Choice) (resolved schema.Meta, err error) {
	errPtr := unsafe.Pointer(&err)
	ident := C.CString(s.WalkState().Position.GetIdent())
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
