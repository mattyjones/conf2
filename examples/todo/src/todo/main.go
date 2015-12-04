// Initialize and start our TODO micro-service application using the Conf2 system
// to load configuration and start management port 

package main

import (
	"flag"
	"data"
	"os"
	"fmt"
	"schema/yang"
	"schema"
	"restconf"
)


// We load from a local config file for simplicity, but same exact data can come
// over network on management port in any accepted format.
var configFileName = flag.String("config", "", "Configuration file")

func main() {
	flag.Parse()

	// TODO: Change this to a file-persistent store.
	if len(*configFileName) == 0 {
		fmt.Fprint(os.Stderr, "Required 'config' parameter missing\n")
		os.Exit(-1)
	}

	// Most applications have a common app service from which you can access
	// all other services and data structures
	app := &App{}
	configFile, err := os.Open(*configFileName)
	if err != nil {
		panic(err)
	}

	// Read json, but you can implement reader in any format you want
	// your reader will be passed schema to validate data.
	config, err := data.NewJsonReader(configFile).Node()

	// load the config into empty app system.  Well designed api will not
	// distinguish config loading from management calls post operation	
	err = data.NodeToNode(config, app.Manage(), app.Schema()).Upsert()

	if err != nil {
		panic(err)
	}
	// start any main thread to keep app from exiting
	app.Start()
}

// Here we mix calls into Conf2 with your application.  You may to separate 
type App struct {
	schema *schema.Module
	management *restconf.Service
	todos *Todos
}

func (app *App) Node() (data.Node) {
	return app.Manage()
}

func (app *App) Schema() schema.MetaList {
	if app.schema == nil {
		var err error
		app.schema, err = yang.LoadModule(yang.YangPath(), "example-todo.yang")
		if err != nil {
			panic(err.Error())
		}
	}
	return app.schema
}

func (app *App) Manage() data.Node {
	s := &data.MyNode{}
	s.OnSelect = func(sel *data.Selection, meta schema.MetaList, new bool) (data.Node, error) {
		switch meta.GetIdent() {
		case "restconf":
			if new {
				app.management = restconf.NewService()
				app.management.RegisterBrowser(app)
			}
			if app.management != nil {
				return app.management.Manage(), nil
			}
		case "todos":
			if new {
				app.todos = &Todos{}
			}
			if app.todos != nil {
				return app.todos.Manage(), nil
			}
		}
		return nil, nil
	}
	return s
}


func (app *App) Start() {
	if app.management == nil {
		panic("Management is not configured")
	}
	app.management.Listen()
}
