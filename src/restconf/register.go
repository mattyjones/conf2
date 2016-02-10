package restconf

import (
	"conf2"
	"data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"schema"
	"strings"
)

// Implements RFC Draft in spirit-only
//   https://tools.ietf.org/html/draft-ietf-netconf-call-home-17
//
// Draft calls for server-initiated registration and this implementation is client-initiated
// which may or may-not be part of the final draft.  Client-initiated registration at first
// glance appears to be more useful, but this may proved to be a wrong assumption on my part.
//
type ControllerRegistry struct {
	Root            data.Data
	CallbackAddress string
	Addr            string
	EndpointId      string
	Registration    *ControllerRegistration
}

type ControllerRegistration struct {
	Id string
}

func (self *ControllerRegistry) Manage() data.Node {
	return &data.Extend{
		Node: data.MarshalContainer(self),
		OnSelect: func(p data.Node, sel *data.Selection, meta schema.MetaList, new bool) (data.Node, error) {
			switch meta.GetIdent() {
			case "registration":
				if self.Registration != nil {
					return data.MarshalContainer(self.Registration), nil
				}
			}
			return nil, nil
		},
		OnEvent: func(p data.Node, sel *data.Selection, e data.Event) error {
			switch e {
			case data.LEAVE_EDIT:
				return self.Register()
			}
			return p.Event(sel, e)
		},
	}
}

func (self *ControllerRegistry) moduleNames(modules schema.MetaList) []string {
	names := make([]string, 0)
	i := schema.NewMetaListIterator(modules, false)
	for i.HasNextMeta() {
		switch mod := i.NextMeta().(type) {
		case *schema.Module:
			names = append(names, mod.GetIdent())
		}
	}
	return names
}

func (self *ControllerRegistry) Register() (err error) {
	var req *http.Request
	conf2.Info.Printf("Registering controller %s", self.Addr)
	if req, err = http.NewRequest("POST", self.Addr, nil); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	modules := strings.Join(self.moduleNames(self.Root.Select().Meta()), ",")
	payload := fmt.Sprintf(`{"modules":["%s"],"endpointId":"%s","callbackAddress":"%s"}`, modules, self.EndpointId,
		self.CallbackAddress)
	req.Body = ioutil.NopCloser(strings.NewReader(payload))
	client := http.DefaultClient
	resp, getErr := client.Do(req)
	if getErr != nil {
		return getErr
	}
	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return conf2.NewErrC(string(respBytes), conf2.Code(resp.StatusCode))
	}
	var rc map[string]interface{}
	if err = json.Unmarshal(respBytes, &rc); err != nil {
		return err
	}
	self.Registration = &ControllerRegistration{
		Id: rc["id"].(string),
	}
	return nil
}
