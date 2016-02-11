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
type CallHome struct {
	Module            *schema.Module
	ControllerAddress string
	EndpointAddress   string
	EndpointId        string
	Registration      *Registration
}

type Registration struct {
	Id string
}

func (self *CallHome) Manage() data.Node {
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
			switch e.Type {
			case data.LEAVE_EDIT:
				return self.Call()
			}
			return p.Event(sel, e)
		},
	}
}

func (self *CallHome) Call() (err error) {
	var req *http.Request
	conf2.Info.Printf("Registering controller %s", self.ControllerAddress)
	if req, err = http.NewRequest("POST", self.ControllerAddress, nil); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	payload := fmt.Sprintf(`{"module":"%s","endpointId":"%s","endpointAddress":"%s"}`, self.Module.GetIdent(),
		self.EndpointId, self.EndpointAddress)
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
	self.Registration = &Registration{
		Id: rc["id"].(string),
	}
	return nil
}
