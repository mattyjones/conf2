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
	CallbackAddress string
	Addr            string
	EndpointId      string
	Registrations   map[string]*ControllerRegistration
}

type ControllerRegistration struct {
	Id string
}

func (self *ControllerRegistry) Manage() data.Node {
	return &data.Extend{
		Node: data.MarshalContainer(self),
		OnSelect: func(p data.Node, sel *data.Selection, meta schema.MetaList, new bool) (data.Node, error) {
			switch meta.GetIdent() {
			case "registrations":
				if self.Registrations != nil {
					return self.ManageRegistrations(), nil
				}
				return nil, nil
			}
			return p.Select(sel, meta, new)
		},
	}
}

func (self *ControllerRegistry) ManageRegistrations() data.Node {
	mm := &data.MarshalMap{Map: self.Registrations}
	return mm.Node()
}

func (self *ControllerRegistry) Register(d data.Data) (err error) {
	var req *http.Request
	conf2.Info.Printf("Registering controller %s", self.Addr)
	if req, err = http.NewRequest("POST", self.Addr, nil); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	module := d.Select().Meta().GetIdent()
	payload := fmt.Sprintf(`{"module":"%s","endpointId":"%s","callbackAddress":"%s"}`, module, self.EndpointId,
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
	reg := &ControllerRegistration{
		Id: rc["id"].(string),
	}
	if self.Registrations == nil {
		self.Registrations = make(map[string]*ControllerRegistration)
	}
	self.Registrations[reg.Id] = reg
	return nil
}
