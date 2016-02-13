package restconf

import (
	"conf2"
	"data"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"schema"
	"time"
	"strings"
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

func NewService(root data.Data) *Service {
	service := &Service{
		Path: "/restconf/",
		Root: root,
		mux:  http.NewServeMux(),
	}
	service.mux.HandleFunc("/.well-known/host-meta", service.resources)
	service.mux.Handle("/restconf/", http.StripPrefix("/restconf/", service))
	service.mux.HandleFunc("/schema/", service.schema)
	return service
}

type Service struct {
	Path            string
	Root            data.Data
	mux             *http.ServeMux
	docrootSource   *docRootImpl
	DocRoot         string
	Port            string
	Iface           string
	CallbackAddress string
	CallHome        *CallHome
}

func (service *Service) EffectiveCallbackAddress() string {
	if len(service.CallbackAddress) > 0 {
		return service.CallbackAddress
	}
	if len(service.Iface) == 0 {
		panic("No iface given for management port")
	}
	ip := conf2.GetIpForIface(service.Iface)
	return fmt.Sprintf("http://%s%s/", ip, service.Port)
}

func (service *Service) Manage() data.Node {
	s := &data.MyNode{Peekables: map[string]interface{}{"internal": service}}
	s.OnSelect = func(sel *data.Selection, r data.ContainerRequest) (data.Node, error) {
		switch r.Meta.GetIdent() {
		case "callHome":
			if r.New {
				service.CallHome = &CallHome{
					EndpointAddress: service.EffectiveCallbackAddress(),
					Module: service.Root.Select().Meta().(*schema.Module),
				}
			}
			if service.CallHome != nil {
				return service.CallHome.Manage(), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *data.Selection, meta schema.HasDataType) (*data.Value, error) {
		return data.ReadField(meta, service)
	}
	s.OnWrite = func(sel *data.Selection, meta schema.HasDataType, v *data.Value) (err error) {
		switch meta.GetIdent() {
		case "docRoot":
			service.DocRoot = v.Str
			service.SetDocRoot(&schema.FileStreamSource{Root: service.DocRoot})
		}
		return data.WriteField(meta, service, v)
	}
	return s
}

type registration struct {
	browser data.Data
}

func (service *Service) handleError(err error, w http.ResponseWriter) {
	if httpErr, ok := err.(data.HttpError); ok {
		http.Error(w, httpErr.Error(), httpErr.HttpCode())
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (service *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var payload data.Node
	var sel data.Selector
	if sel = service.Root.Select().Selector().FindUrl(r.URL); sel.LastErr == nil {
		if sel.Selection == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err != nil {
			service.handleError(err, w)
			return
		}
		switch r.Method {
		case "DELETE":
			err = sel.Selection.Delete()
		case "GET":
			w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
			output := data.NewJsonWriter(w).Node()
			err = sel.Push(output).Insert().LastErr
		case "PUT":
			err = sel.Pull(data.NewJsonReader(r.Body).Node()).Upsert().LastErr
		case "POST":
			if schema.IsAction(sel.Selection.Meta()) {
				input := data.NewJsonReader(r.Body).Node()
				var outputSel *data.Selection
				outputSel, err = sel.Selection.Action(input)
				if outputSel != nil {
					w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
					err = outputSel.Selector().Push(data.NewJsonWriter(w).Node()).Insert().LastErr
				}
			} else {
				payload = data.NewJsonReader(r.Body).Node()
				err = sel.Pull(payload).Insert().LastErr
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else {
		err = sel.LastErr
	}

	if err != nil {
		service.handleError(err, w)
	}
}

type docRootImpl struct {
	docroot schema.StreamSource
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

func (service *Service) schema(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	if p := strings.TrimPrefix(r.URL.Path, "/schema/"); len(p) < len(r.URL.Path) {
		r.URL.Path = p
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	m := service.Root.Select().Meta().(*schema.Module)
	sch := data.NewSchemaData(m, false)
	if sel := sch.Select().Selector().FindUrl(r.URL); sel.LastErr != nil {
		service.handleError(sel.LastErr, w)
		return
	} else if sel.Selection == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
		output := data.NewJsonWriter(w).Node()
		if sel = sel.Push(output).Insert(); sel.LastErr != nil {
			service.handleError(sel.LastErr, w)
			return
		}
	}
}

func (service *Service) resources(w http.ResponseWriter, r *http.Request) {
	// RESTCONF Sec. 3.1
	fmt.Fprintf(w, `"xrd" : { "link" : { "@rel" : "restconf", "@href" : "/restconf" } } }`)
}
