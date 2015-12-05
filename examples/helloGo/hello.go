// This is a very simple example of using the Conf2 library to add RESTCONF
// support to an app that can say hello to you.
//
// To run this example, use:
//    GOPATH=$(readlink -f ../../) go run hello.go
//

package main

import (
	"restconf"
	"data"
	"schema"
	"schema/yang"
)

// This is YANG and usually stored in files but you can store your yang anywhere
// you want including in your source or in a database.  See YANG RFC for full
// options of definitions including typedefs, groupings, containers, lists
// actions, leaf-lists, choices, enumerations and others.
var helloApiDefinition = `
/*
  module is a collection or definitions, much like a YANG container except its the
  top-most container
*/
module hello {
  prefix "hello";
  namespace "hello";
  revision 0;

  /* a "leaf" is just a field */
  leaf message {
  	type string;
  }

  leaf count {
    /* marking fields as NOT config helps denotes them as metrics */
    config "false";
    type int32;
  }

  /* An "rpc" is a function you want to expose */
  rpc say {

    /* rpcs can have optional input defined */
    input {
      leaf name {
        type string;
      }
    }

    /* rpcs can have optional output defined */
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

	// Your app, no references to Conf2 enc.
	app := &MyApp{}

	// This is the connection between your app and Conf2.  ManageApp can then
	// navigate through your app's other structures to fulfil API.
	manage := &ManageMyApp{App:app}

	// Here we load from memory, but to load from YANGPATH environment variable use:
	//  yang.LoadModule("some-module", data.YangPath())
	if manage.Meta, err = yang.LoadModuleFromByteArray([]byte(helloApiDefinition), nil); err != nil {
		panic(err.Error())
	}

	// You can register as many APIs as you want, The module name is the default RESTCONF base url
	if err = service.RegisterBrowser(manage); err != nil {
		panic(err.Error())
	}

	// you may want to start in background, but here we start in foreground to keep app running.
	// Hit Ctrl-C in terminal to quit
	service.Listen()
}

// Beginning of your existing application code and has no references to Conf2
type MyApp struct {
	Message string
	Count int
}

// a random function we'll expose thru API using OnAction below
func (app *MyApp) SayHello(name string) string {

	// as a random metric, let's count how many times we've said hello
	app.Count++

	return app.Message + " " + name
}

// Beginning of Conf2 code
// This Go struct implements data.Data interface, the entry point for your API
type ManageMyApp struct {
	App *MyApp
	Meta *schema.Module
}

// This implements data.Data get's the initial node into your app
func (manage *ManageMyApp) Node() (data.Node) {
	return manage.Manage()
}

// This implements gets your app's meta data loaded from your yang definition
func (manage *ManageMyApp) Schema() schema.MetaList {
	return manage.Meta
}

// Each unique Go struct you want to manage typically has a corresponding management Node
// but similar structs may share Node definitions.
func (manage *ManageMyApp) Manage() data.Node {

	// Node is an interface, but there's a convenient struct that implements the interface
	// and delegates operations to closure functions defined below
	n := &data.MyNode{}

	// If we had nested Go structs, we'd implement OnSelect to drill into data
	// If we were managing a list of structs, we'd implement OnNext

	// This is for reading data, URL is
	//
	//   GET http://localhost:8009/restconf/hello
	//   Example Response:
	//    {"message":"hello","count":10}
	//
	n.OnRead = func(sel *data.Selection, meta schema.HasDataType) (*schema.Value, error) {
		// Here we can use Go's reflection but if reflection isn't valid for some or all fields,
		// you can add a switch case to handle them separately
		return data.ReadField(meta, manage.App)
	}

	// This is for reading data, URL is
	//
	//   PUT http://localhost:8009/restconf/hello
	//
	//   Example Payload:
	//     {"message":"hello"}
	//
	n.OnWrite = func(sel *data.Selection, meta schema.HasDataType, v *schema.Value) error {
		// Here we can use Go's reflection but if names are different, add a switch case
		// statement here
		return data.WriteField(meta, manage.App, v)
	}

	// Any RPCs (YANG "rpc" or "action") with come thru here
	//
	//  POST http://localhost:8009/restconf/hello/say
	//
	//   Example Payload:
	//     {"name":"joe"}
	///
	//   Example Response:
	//     {"message":"hello joe"}
	//
	n.OnAction = func(sel *data.Selection, rpc *schema.Rpc, input data.Node) (output data.Node, err error) {

		// You can use a variety of methods to unmarshal the input including sticking into go map
		// using the Bucket struct
		param := struct { Name string } {}
		err = data.NodeToNode(input, data.MarshalContainer(input), rpc.Input).Insert()

		// See how we can call functions of our app from our management data
		// browser?
		s := manage.App.SayHello(param.Name)

		// Build the response, we choose reflection marshaller again just like input data
		response := struct { Message string } { s }
		return data.MarshalContainer(&response), nil
	}

	return n
}
