package main
import (
	"restconf"
	"data"
	"schema"
	"schema/yang"
)

// This is YANG and usually stored in files but you can store your yang anywhere
// you want including in your source or in a database
var api = `
module hello {
  prefix "hello";
  namespace "hello";
  revision 0;

  leaf message {
  	type string;
  }

  leaf count {
    config "false";
    type int32;
  }

  rpc say {
    input {
      leaf name {
        type string;
      }
    }
    output {
      leaf message {
        type string;
      }
    }
  }
}
`

func main() {
	var err error

	// Create a RESTCONF service to register your APIs
	service := restconf.NewService()
	service.Port = ":8009"

	app := &MyApp{}
	manage := &ManageMyApp{App:app}
	if manage.Meta, err = yang.LoadModuleFromByteArray([]byte(api), nil); err != nil {
		panic(err.Error())
	}
	if err = service.RegisterBrowser(manage); err != nil {
		panic(err.Error())
	}

	// you may want to start in background, but here we start in foreground to keep app running.
	// Hit Ctrl-C in terminal to quit
	service.Listen()
}

// Your existing application code void of any Conf2 code
// ============================================================
type MyApp struct {
	Message string
	Count int
}

func (app *MyApp) SayHello(name string) string {
	app.Count++
	return app.Message + " " + name
}

// Begin Conf2 code

type ManageMyApp struct {
	App *MyApp
	Meta *schema.Module
}

func (manage *ManageMyApp) Selector(path *data.Path) (*data.Selection, error) {
	return data.WalkPath(data.NewSelection(manage.Manage(), manage.Meta), path)
}

func (manage *ManageMyApp) Schema() schema.MetaList {
	return manage.Meta
}

// Each unique Go struct you want to manage typically has a corresponding management Node
// but similar structs may share Node definitions.
func (manage *ManageMyApp) Manage() data.Node {

	// Node is an interface, but there's a convenient struct that implements the interface
	// and delegates operations to anonymous functions defined within
	n := &data.MyNode{}

	// If we had nested Go structs, we'd implement OnSelect to drill into data
	// If we were managing a list of structs, we'd implement OnNext

	n.OnRead = func(sel *data.Selection, meta schema.HasDataType) (*data.Value, error) {
		// Here we can use Go's reflection but if names are different, add a switch case
		// statement here
		return data.ReadField(meta, manage.App)
	}

	n.OnWrite = func(sel *data.Selection, meta schema.HasDataType, v *data.Value) error {
		// Here we can use Go's reflection but if names are different, add a switch case
		// statement here
		return data.WriteField(meta, manage.App, v)
	}

	// Any RPCs (YANG "rpc" or "action") with come thru here
	n.OnAction = func(sel *data.Selection, rpc *schema.Rpc, input *data.Selection) (output *data.Selection, err error) {
		param := struct {
			Name string
		} {}
		err = data.MarshalTo(input, &param)
		manage.App.SayHello(param.Name)
		response := struct {
			Message string
		} {
			manage.App.SayHello(param.Name),
		}
		return data.NewSelection(data.MarshalContainer(&response), rpc.Output), nil
	}

	return n
}
