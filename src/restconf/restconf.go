package restconf

import (
	"schema"
	"schema/browse"
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
	RegisterBrowserWithName(browser browse.Browser, name string) error
	SetDocRoot(schema.StreamSource)
	Stop()
}

func NewService() (Service, error) {
	service := &serviceImpl{restconfPath:"/restconf/"}
	service.registrations = make(map[string]*registration, 5)
	service.mux = http.NewServeMux()
	service.mux.HandleFunc("/.well-known/host-meta", service.resources)
	// always add browser for restconf server itself
	rcb, err := NewBrowser(service)
	if err != nil {
		return nil, err
	}
	if err = service.RegisterBrowser(rcb); err != nil {
		return nil, err
	}
	return service, nil
}

type serviceImpl struct {
	restconfPath string
	registrations map[string]*registration
	mux *http.ServeMux
	docroot *docRootImpl
}

type registration struct {
	browser browse.Browser
}

func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var path *browse.Path
	if path, err = browse.ParsePath(r.URL.Path); err == nil {
		var selection browse.Selection
		var state *browse.WalkState
		if selection, state, err = reg.browser.RootSelector(); err == nil {
			if selection, state, err = browse.WalkPath(state, selection, path); err == nil {
				var walkCntlr browse.WalkController
				if selection == nil {
					http.Error(w, r.URL.Path, http.StatusNotFound)
				} else {
					switch r.Method {
					case "GET":
						w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
						wtr := browse.NewJsonWriter(w)
						var out browse.Selection
						if out, err = wtr.GetSelector(); err == nil {
							walkCntlr = browse.NewFullWalk(r.URL.RawQuery)
							err = browse.Upsert(state, selection, out, walkCntlr)
						}
					case "POST", "PUT":
						rdr := browse.NewJsonReader(r.Body)
						var in browse.Selection
						if in, err = rdr.GetSelector(state); err == nil {
							walkCntlr = browse.WalkAll()
							switch r.Method {
							case "POST":
								err = browse.Insert(state, in, selection, walkCntlr)
							case "PUT":
								err = browse.Upsert(state, in, selection, walkCntlr)
							}

							if err == nil {
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
	docroot schema.StreamSource
}

func (service *serviceImpl) RegisterBrowser(browser browse.Browser) error {
	ident := browser.Module().GetIdent()
	return service.RegisterBrowserWithName(browser, ident)
}

func (service *serviceImpl) RegisterBrowserWithName(browser browse.Browser, ident string) error {
	reg := &registration{browser}
	service.registrations[ident] = reg
	fullPath := fmt.Sprint(service.restconfPath, ident, "/")
	log.Println("registering browser at path ", fullPath)
	service.mux.Handle(fullPath,  http.StripPrefix(fullPath, reg))
	return nil
}

func (service *serviceImpl) SetDocRoot(docroot schema.StreamSource) {
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
	log.Println("Starting RESTCONF interface")
	log.Fatal(s.ListenAndServe())
}

func (service *serviceImpl) Stop() {
	if service.docroot != nil && service.docroot.docroot != nil {
		schema.CloseResource(service.docroot.docroot)
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
		defer schema.CloseResource(rdr)
		ext := filepath.Ext(path)
		ctype := mime.TypeByExtension(ext)
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
