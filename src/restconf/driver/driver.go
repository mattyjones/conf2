package driver

import "C"

import (
	"restconf"
	"yang/driver"
	"yang"
	"unsafe"
	"fmt"
	"yang/browse"
)

//export restconfc2_service_new
func restconfc2_service_new() unsafe.Pointer {
	service, err := restconf.NewService()
	// TODO : Proper error handling
	if err != nil {
		return nil
	}
	return driver.NewGoHandle(service).ID
}

//export restconfc2_service_start
func restconfc2_service_start(service_hnd_id unsafe.Pointer) {
	service_hnd, found := driver.GoHandles()[service_hnd_id]
	if ! found {
		panic(fmt.Sprint("Restconf service not found", service_hnd_id))
	}
	service := service_hnd.Data.(restconf.Service)
	go service.Listen()
}

//export restconfc2_set_doc_root
func restconfc2_set_doc_root(service_hnd_id unsafe.Pointer, stream_source_hnd_id unsafe.Pointer) {
	service_hnd, found := driver.GoHandles()[service_hnd_id]
	if ! found {
		panic(fmt.Sprint("Restconf service not found", service_hnd_id))
	}
	service := service_hnd.Data.(restconf.Service)

	stream_source_hnd, found := driver.GoHandles()[stream_source_hnd_id]
	if ! found {
		panic(fmt.Sprint("Restconf service not found", service_hnd_id))
	}
	stream_source := stream_source_hnd.Data.(yang.StreamSource)
	service.SetDocRoot(stream_source)
}

//export restconfc2_register_browser
func restconfc2_register_browser(service_hnd_id unsafe.Pointer, browser_hnd_id unsafe.Pointer) error {
	service_hnd, found := driver.GoHandles()[service_hnd_id]
	if ! found {
		panic(fmt.Sprint("Restconf service not found", service_hnd_id))
	}
	service := service_hnd.Data.(restconf.Service)

	browser_hnd, found := driver.GoHandles()[browser_hnd_id]
	if ! found {
		panic(fmt.Sprint("Restconf service not found", browser_hnd_id))
	}
	browser := browser_hnd.Data.(browse.Browser)

	return service.RegisterBrowser(browser)
}