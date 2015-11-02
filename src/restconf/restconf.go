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

func NewService() (*Service, error) {
	service := &Service{restconfPath:"/restconf/"}
	service.registrations = make(map[string]*registration, 5)
	service.mux = http.NewServeMux()
	service.mux.HandleFunc("/.well-known/host-meta", service.resources)
	// always add browser for restconf server itself
	rcb, err := NewData(service)
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
	browser browse.Data
}


func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleError := func (err error) {
		if httpErr, ok := err.(browse.HttpError); ok {
			http.Error(w, httpErr.Error(), httpErr.HttpCode())
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	var err error
	var path *browse.Path
	if path, err = browse.ParsePath(r.URL.Path); err == nil {
		var selection *browse.Selection
		if selection, err = reg.browser.Selector(path); err != nil {
			handleError(err)
			return
		}
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
			output := selection.Copy(browse.NewJsonWriter(w).Container())
			err = browse.ControlledInsert(selection, output, browse.LimitedWalk(r.URL.RawQuery))
		case "PUT":{
			var payload browse.Node
			if payload, err = browse.NewJsonReader(r.Body).NodeFromSelection(selection); err != nil {
				handleError(err)
				return
			}
			err = browse.UpsertByNode(selection, payload, selection.Node())
		}
		case "POST": {
			if schema.IsAction(selection.SelectedMeta()) {
				var rpcInput browse.Node
				var rpcOutput *browse.Selection
				if rpcInput, err = browse.NewJsonReader(r.Body).Node(); err != nil {
					handleError(err)
					return
				}
				if rpcOutput, err = browse.Action(selection, rpcInput); err != nil {
					handleError(err)
					return
				}
				if rpcOutput != nil {
					w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
					output := rpcOutput.Copy(browse.NewJsonWriter(w).Container())
					browse.Insert(rpcOutput, output)
				}
			} else {
				var payload browse.Node
				if payload, err = browse.NewJsonReader(r.Body).NodeFromSelection(selection); err != nil {
					handleError(err)
					return
				}
				err = browse.InsertByNode(selection, payload, selection.Node())
			}
		}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}

	if err != nil  {
		handleError(err)
	}
}

type docRootImpl struct {
	docroot schema.StreamSource
}

func (service *Service) RegisterBrowser(browser browse.Data) error {
	ident := browser.Schema().GetIdent()
	return service.RegisterBrowserWithName(browser, ident)
}

func (service *Service) RegisterBrowserWithName(browser browse.Data, ident string) error {
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
