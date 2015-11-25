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

var configFileName = flag.String("config", "", "Configuration file")

func main() {
	flag.Parse()

	// TODO: Change this to a file-persistent store.
	if len(*configFileName) == 0 {
		fmt.Fprint(os.Stderr, "Required 'config' parameter missing\n")
		os.Exit(-1)
	}

	app := &App{}
	selection, err := app.Selector(data.NewPath(""))
	if err != nil {
		panic(err)
	}

	configFile, err := os.Open(*configFileName)
	if err != nil {
		panic(err)
	}
	config, err := data.NewJsonReader(configFile).Selector(selection.State)
	err = data.Insert(config, selection)
	if err != nil {
		panic(err)
	}

	app.Start()
}

type App struct {
	schema *schema.Module
	management *restconf.Service
	todos *Todos
}

func (app *App) Selector(path *data.Path) (*data.Selection, error) {
	return data.WalkPath(data.NewSelection(app.Manage(), app.Schema()), path)
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