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
	Code int
	Msg string
}

func (err *restconfError) Error() string {
	return err.Msg
}

func (err *restconfError) HttpCode() int {
	return err.Code
}

//type Service interface {
//	Listen()
//	RegisterBrowser(browser browse.Browser) error
//	RegisterBrowserWithName(browser browse.Browser, name string) error
//	SetDocRoot(schema.StreamSource)
//	Stop()
//}

func NewService() (*Service, error) {
	service := &Service{restconfPath:"/restconf/"}
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

type Service struct {
	restconfPath string
	registrations map[string]*registration
	mux *http.ServeMux
	docroot *docRootImpl
	Port string
}

type registration struct {
	browser browse.Browser
}

func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var path *browse.Path
	if path, err = browse.ParsePath(r.URL.Path); err == nil {
		path.SetQuery(r.URL.RawQuery)
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
			dest := browse.NewJsonFragmentWriter(w)
			err = browse.Insert(path, reg.browser, dest)
		case "PUT":{
			rdr := browse.NewJsonFragmentReader(r.Body)
			err = browse.Upsert(path, rdr, reg.browser)
		}
		case "POST": {
			rdr := browse.NewJsonFragmentReader(r.Body)
			dest := browse.NewJsonFragmentWriter(w)
			err = reg.operation(path, rdr, reg.browser, dest)
		}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	if err != nil  {
		if httpErr, ok := err.(browse.HttpError); ok {
			http.Error(w, httpErr.Error(), httpErr.HttpCode())
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (reg *registration) operation(path *browse.Path, src browse.Browser, dest browse.Browser, output browse.Browser) (err error) {
	var destSel, srcSel browse.Selection
	var destState, srcState *browse.WalkState
	if destSel, destState, err = dest.Selector(path, browse.INSERT); err != nil {
		return err
	}
	if destSel == nil {
		return browse.NotFound(path.URL)
	}
	state := destState
	if state == nil {
		state = srcState
		if state == nil {
			return browse.NotFound(path.URL)
		}
	}
	if srcSel, srcState, err = src.Selector(path, browse.READ); err != nil {
		return err
	}
	if destState.Position() != nil && schema.IsAction(destState.Position()) {
		var outputSel, rpcOutput browse.Selection
		var outputState *browse.WalkState
		actionMeta := state.Position().(*schema.Rpc)
		if rpcOutput, outputState, err = destSel.Action(state, actionMeta, srcSel); err != nil {
			return err
		}
		if rpcOutput != nil {
			outputSel, _, err = output.Selector(browse.NewPath(""), browse.INSERT)
			if err = browse.Edit(outputState, rpcOutput, outputSel, browse.INSERT, browse.LimitedWalk(path.Query)); err != nil {
				return err
			}
		}
	} else {
		err = browse.Edit(state, srcSel, destSel, browse.INSERT, browse.LimitedWalk(path.Query))
	}
	return
}

type docRootImpl struct {
	docroot schema.StreamSource
}

func (service *Service) RegisterBrowser(browser browse.Browser) error {
	ident := browser.Schema().GetIdent()
	return service.RegisterBrowserWithName(browser, ident)
}

func (service *Service) RegisterBrowserWithName(browser browse.Browser, ident string) error {
	reg := &registration{browser}
	service.registrations[ident] = reg
	fullPath := fmt.Sprint(service.restconfPath, ident, "/")
	log.Println("registering browser at path ", fullPath)
	service.mux.Handle(fullPath,  http.StripPrefix(fullPath, reg))
	return nil
}

func (service *Service) SetDocRoot(docroot schema.StreamSource) {
	service.docroot = &docRootImpl{docroot:docroot}
	service.mux.Handle("/ui/", http.StripPrefix("/ui/", service.docroot))
}

func (service *Service) AddHandler(pattern string, handler http.Handler) {
	service.mux.Handle(pattern, http.StripPrefix(pattern, handler))
}

func (service *Service) Listen() {
	s := &http.Server{
		Addr:           service.Port,
		Handler:        service.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Starting RESTCONF interface")
	log.Fatal(s.ListenAndServe())
}

func (service *Service) Stop() {
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

func (service *Service) resources(w http.ResponseWriter, r *http.Request) {
	// RESTCONF Sec. 3.1
	fmt.Fprintf(w, `"xrd" : { "link" : { "@rel" : "restconf", "@href" : "/restconf" } } }`)
}
