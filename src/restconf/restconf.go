package restconf

import (
	"conf2"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"schema"
	"schema/browse"
	"time"
)

type restconfError struct {
	Code int
	Msg  string
}

func (err *restconfError) Error() string {
	return err.Msg
}

func (err *restconfError) HttpCode() int {
	return err.Code
}

func NewService() *Service {
	service := &Service{Path: "/restconf/"}
	service.registrations = make(map[string]*registration, 5)
	service.mux = http.NewServeMux()
	service.mux.HandleFunc("/.well-known/host-meta", service.resources)
	return service
}

type Service struct {
	Path          string
	registrations map[string]*registration
	mux           *http.ServeMux
	docrootSource *docRootImpl
	DocRoot		  string
	Port          string
}

func (service *Service) Manage() browse.Node {
	s := &browse.MyNode{}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "registrations":
			strlist := make([]string, len(service.registrations))
			i := 0
			for name, _ := range service.registrations {
				strlist[i] = name
				i++
			}
			return &browse.Value{Strlist:strlist}, nil
		default:
			return browse.ReadField(meta, service)
		}
		return nil, nil
	}
	s.OnWrite = func(sel *browse.Selection, meta schema.HasDataType, v *browse.Value) (err error) {
		switch meta.GetIdent() {
		case "docRoot":
			service.DocRoot = v.Str
			service.SetDocRoot(&schema.FileStreamSource{Root:service.DocRoot})
		}
		return browse.WriteField(meta, service, v)
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) (err error) {
		switch e {
		case browse.NEW:
			rcb, err := NewData(service)
			if err != nil {
				return err
			}
			// always add browser for restconf server itself
			if err = service.RegisterBrowser(rcb); err != nil {
				return err
			}
		}
		return
	}

	return s
}

type registration struct {
	browser browse.Data
}

func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleError := func(err error) {
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
		case "DELETE":
			err = browse.Delete(selection, selection.Node())
		case "GET":
			w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
			output := selection.Copy(browse.NewJsonWriter(w).Container())
			err = browse.ControlledInsert(selection, output, browse.LimitedWalk(r.URL.RawQuery))
		case "PUT":
			{
				var payload browse.Node
				if payload, err = browse.NewJsonReader(r.Body).NodeFromSelection(selection); err != nil {
					handleError(err)
					return
				}
				err = browse.UpsertByNode(selection, payload, selection.Node())
			}
		case "POST":
			{
				if schema.IsAction(selection.Position()) {
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

	if err != nil {
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
	fullPath := fmt.Sprint(service.Path, ident, "/")
	conf2.Info.Println("registering browser at path ", fullPath)
	service.mux.Handle(fullPath, http.StripPrefix(fullPath, reg))
	return nil
}

func (service *Service) SetDocRoot(docroot schema.StreamSource) {
	service.docrootSource = &docRootImpl{docroot: docroot}
	service.mux.Handle("/ui/", http.StripPrefix("/ui/", service.docrootSource))
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
	conf2.Info.Println("Starting RESTCONF interface")
	conf2.Err.Fatal(s.ListenAndServe())
}

func (service *Service) Stop() {
	if service.docrootSource != nil && service.docrootSource.docroot != nil {
		schema.CloseResource(service.docrootSource.docroot)
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
