package c

// #include "conf2/handle.h"
// extern void conf2_handle_release_bridge(conf2_handle_release_impl impl_func, void *handle, void *errPtr);
import "C"

import (
	"unsafe"
	"fmt"
	"schema"
	"sync"
)

type Handle interface {
	schema.Resource
}

var NilHandle unsafe.Pointer

func NewGoHandle(data interface{}) *GoHandle {
	hnd := &GoHandle{ID:unsafe.Pointer(&data), Data:data}
	GoHandles()[hnd.ID] = hnd
	return hnd
}

type GoHandle struct {
	ID unsafe.Pointer
	Data interface{}
}

func (hnd *GoHandle) Close() error {
	// just removing a reference from it allows the GC to do the rest
	delete(GoHandles(), hnd.ID)
	return nil
}

type ApiHandle struct {
	ID unsafe.Pointer
	release_impl C.conf2_handle_release_impl
}

func (hnd *ApiHandle) Close() (err error) {
	if hnd.ID == nil {
		panic("Attempting to Close unitialized or pre-Closed handle")
	}
	if hnd.release_impl != nil {
		C.conf2_handle_release_bridge(hnd.release_impl, hnd.ID, unsafe.Pointer(&err))
	}
	delete(apiHandles, hnd.ID)
	hnd.ID = nil
	return err
}

var apiHandles map[unsafe.Pointer]*ApiHandle
var goHandles map[unsafe.Pointer]*GoHandle

var initGoHandles sync.Once

func ApiHandles() map[unsafe.Pointer]*ApiHandle {
	if apiHandles == nil {
		apiHandles = make(map[unsafe.Pointer]*ApiHandle, 100)
	}
	return apiHandles
}

func GoHandles() map[unsafe.Pointer]*GoHandle {
	if goHandles == nil {
		initGoHandles.Do(func() {
			goHandles = make(map[unsafe.Pointer]*GoHandle, 100)
		})
	}
	return goHandles
}

//export conf2_handle_new
func conf2_handle_new(api_handle unsafe.Pointer, release_impl C.conf2_handle_release_impl) unsafe.Pointer {
	go_handle := &ApiHandle{ID:api_handle, release_impl:release_impl}
	key := unsafe.Pointer(go_handle)
	ApiHandles()[key] = go_handle
	return key
}

//export conf2_handle_release
func conf2_handle_release(key unsafe.Pointer) {
	handle, valid := GoHandles()[key]
	if valid {
		handle.Close()
	} else {
		panic(fmt.Sprint("Close on invalid handle", key))
	}
}
