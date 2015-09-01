package restconf

import (
	"yang"
	"yang/browse"
	"net/http"
	"time"
	"fmt"
	"log"
	"io"
	"mime"
	"path/filepath"
)

type restconfError struct {
	Code browse.ResponseCode
	Msg string
}

func (err *restconfError) Error() string {
	return err.Msg
}

type Service interface {
	Listen()
	RegisterBrowser(browser browse.Browser) error
	SetDocRoot(yang.StreamSource)
	Stop()
}

func NewService() Service {
	service := &serviceImpl{restconfPath:"/restconf/"}
	service.registrations = make(map[string]registration, 5)
	service.mux = http.NewServeMux()
	service.mux.HandleFunc("/.well-known/host-meta", service.resources)
	return service
}

type serviceImpl struct {
	restconfPath string
	registrations map[string]registration
	mux *http.ServeMux
	docroot *docRootImpl
}

type registration struct {
	browser browse.Browser
}

func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTCONF", r.URL.Path)
	var err error
	if path, err := browse.NewPath(r.URL.Path); err == nil {
		if selection, err := reg.browser.RootSelector(); err == nil {
			if selection, err = browse.WalkPath(selection, path); err == nil {
				if selection == nil {
					http.Error(w, r.URL.Path, http.StatusNotFound)
				} else {
					switch r.Method {
					case "GET":
						w.Header().Set("Content-Type", mime.TypeByExtension("json"))
						wtr := browse.NewJsonWriter(w)
						if out, err := wtr.GetSelector(); err == nil {
							err = browse.Insert(selection, out)
						}
					case "POST":
						rdr := browse.NewJsonReader(r.Body)
						if in, err := rdr.GetSelector(selection.Meta); err == nil {
							if err = browse.Insert(in, selection); err == nil {
								http.Error(w, "", http.StatusNoContent)
							}
						}
					default:
						http.Error(w, "Not implemented yet", http.StatusInternalServerError)
					}
				}
			}
		}
	}
	if err != nil  {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type docRootImpl struct {
	docroot yang.StreamSource
}

func (service *serviceImpl) RegisterBrowser(browser browse.Browser) error {
	reg := registration{browser}
	ident := browser.Module().GetIdent()
	service.registrations[ident] = reg
	fullPath := fmt.Sprint(service.restconfPath, ident, "/")
	log.Println("registering browser at path ", fullPath)
	service.mux.Handle(fullPath,  http.StripPrefix(fullPath, &reg))
	return nil
}

func (service *serviceImpl) SetDocRoot(docroot yang.StreamSource) {
	service.docroot = &docRootImpl{docroot:docroot}
	service.mux.Handle("/ui/", http.StripPrefix("/ui/", service.docroot))
}

func (service *serviceImpl) Listen() {
	s := &http.Server{
		Addr:           ":8008",
		Handler:        service.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func (service *serviceImpl) Stop() {
	if service.docroot != nil && service.docroot.docroot != nil {
		service.docroot.docroot.Close()
	}
	// TODO - actually stop service
}

func (service *docRootImpl) ServeHTTP(wtr http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path == "" {
		path = "index.html"
	}
	if rdr, err := service.docroot.OpenStream(path); err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	} else {
		defer rdr.Close()
		ctype := mime.TypeByExtension(filepath.Ext(path))
		wtr.Header().Set("Content-Type", ctype)
		if _, err = io.Copy(wtr, rdr); err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
		}
		// Eventually support this but need file seeker to do that.
		// http.ServeContent(wtr, req, path, time.Now(), &ReaderPeeker{rdr})
	}
}

func (service *serviceImpl) resources(w http.ResponseWriter, r *http.Request) {
	// RESTCONF Sec. 3.1
	fmt.Fprintf(w, `"xrd" : { "link" : { "@rel" : "restconf", "@href" : "/restconf" } } }`)
}
