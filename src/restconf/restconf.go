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
	Path            string
	registrations   map[string]*registration
	mux             *http.ServeMux
	docrootSource   *docRootImpl
	DocRoot         string
	Port            string
	Iface           string
	CallbackAddress string
	Controller      *ControllerRegistry
}

func (service *Service) EffectiveCallbackAddress() string {
	if len(service.CallbackAddress) > 0 {
		return service.CallbackAddress
	}
	if len(service.Iface) == 0 {
		panic("No iface given for management port")
	}
	ip := conf2.GetIpForIface(service.Iface)
	return fmt.Sprintf("http://%s%s/restconf", ip, service.Port)
}

func (service *Service) Manage() data.Node {
	s := &data.MyNode{}
	s.OnSelect = func(sel *data.Selection, meta schema.MetaList, new bool) (data.Node, error) {
		switch meta.GetIdent() {
		case "controller":
			if new {
				service.Controller = &ControllerRegistry{
					CallbackAddress: service.EffectiveCallbackAddress(),
				}
			}
			if service.Controller != nil {
				return service.Controller.Manage(), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *data.Selection, meta schema.HasDataType) (*data.Value, error) {
		switch meta.GetIdent() {
		case "registrations":
			strlist := make([]string, len(service.registrations))
			i := 0
			for name, _ := range service.registrations {
				strlist[i] = name
				i++
			}
			return &data.Value{Strlist: strlist}, nil
		default:
			return data.ReadField(meta, service)
		}
		return nil, nil
	}
	s.OnWrite = func(sel *data.Selection, meta schema.HasDataType, v *data.Value) (err error) {
		switch meta.GetIdent() {
		case "docRoot":
			service.DocRoot = v.Str
			service.SetDocRoot(&schema.FileStreamSource{Root: service.DocRoot})
		}
		return data.WriteField(meta, service, v)
	}
	s.OnEvent = func(sel *data.Selection, e data.Event) (err error) {
		switch e {
		case data.NEW:
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
	browser data.Data
}

func (reg *registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleError := func(err error) {
		if httpErr, ok := err.(data.HttpError); ok {
			http.Error(w, httpErr.Error(), httpErr.HttpCode())
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	var err error
	var payload data.Node
	sel := reg.browser.Select()
	if sel, err = sel.FindUrl(r.URL); err == nil {
		if sel == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err != nil {
			handleError(err)
			return
		}
		switch r.Method {
		case "DELETE":
			err = sel.Delete()
		case "GET":
			w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
			output := data.NewJsonWriter(w).Node()
			err = sel.Push(output).ControlledInsert(data.LimitedWalk(r.URL.Query()))
		case "PUT":
			err = sel.Pull(data.NewJsonReader(r.Body).Node()).Upsert()
		case "POST":
			if schema.IsAction(sel.Meta()) {
				input := data.NewJsonReader(r.Body).Node()
				var outputSel *data.Selection
				outputSel, err = sel.Action(input)
				if outputSel != nil {
					w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
					err = outputSel.Push(data.NewJsonWriter(w).Node()).Insert()
				}
			} else {
				payload = data.NewJsonReader(r.Body).Node()
				err = sel.Pull(payload).Insert()
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

func (service *Service) RegisterBrowser(browser data.Data) error {
	ident := browser.Select().Meta().GetIdent()
	return service.RegisterBrowserWithName(browser, ident)
}

func (service *Service) RegisterBrowserWithName(browser data.Data, ident string) error {
	reg := &registration{browser}
	service.registrations[ident] = reg
	fullPath := fmt.Sprint(service.Path, ident, "/")
	conf2.Info.Println("registering browser at path ", fullPath)
	service.mux.Handle(fullPath, http.StripPrefix(fullPath, reg))
	if service.Controller != nil {
		// TODO: Controller register should have background service that refreshes
		// registration on it's own.  Do it once for now.
		go func() {
			if remoteRegErr := service.Controller.Register(browser); remoteRegErr != nil {
				conf2.Err.Fatal(remoteRegErr)
			}
		}()
	}
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
