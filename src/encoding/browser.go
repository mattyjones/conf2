package oasm

// #include "encoding/c2_encoding_visit.h"
// extern int c2_encoding_visit(c2_encoding_visit_impl);
import "C"

import (
	"yang"
	"io"
	"unsafe"
)

type Selector interface {
	Selected(schema yang.Def, data interface{})
}

type Path struct {
	parent *Path
}

type Visitor interface {
	EnterContainer(*yang.Container, interface{}) bool
	EnterList(*yang.List, []interface{}) bool
	VisitIntLeaf(*yang.Leaf, int)
	VisitStringLeaf(*yang.Leaf, string)
	VisitStringLeafList(*yang.Leaf, []string)
	VisitIntLeafList(*yang.Leaf, []int)
	LeaveContainer(*yang.Container, interface{})
	LeaveList(*yang.List, interface{})
}

type JsonBuilder struct {
	schema *yang.Module
	wtr io.Writer
}

type DriverVisitor struct {
	impl C.c2_encoding_visit_impl
	obj unsafe.Pointer
}

func (drv *DriverVisitor) VisitContainer(p Path, ident string, container Browser) {

}

type UriSelector struct {
	uri string
	visitor Visitor
}

func (s *UriSelector) Enter(node interface{}) {

}

type Browser interface {
	Accept(visitor Visitor)
}

type SchemaBrowser struct {
	module *yang.Module
	schema *yang.Module
}

func (browser *SchemaBrowser) AcceptDefinitionHeader(visitor Visitor, selector Selector, schema yang.Def, def yang.Def) {
	switch schema.GetIdent() {
		case "ident":
			visitor.VisitStringLeaf(schema, def.GetIdent())
		case "description":
			visitor.VisitStringLeaf(schema, def.(yang.Describable).GetDescription())
		default :
			// FAIL
	}
}

func (browser *SchemaBrowser) AcceptDataDefinition(visitor Visitor, selector Selector, schema yang.Def, def yang.DefList) {
	dataDef := def.(yang.DefList).GetFirstDef()
	for dataDef != nil {
		choice := schema.(yang.DefList).GetFirstDef().(*yang.Choice) // "body-stmt"
		switch t := dataDef.(type) {
		default:
		case *yang.Container:
			// this really should be Yang:Choice/Case
			childSchema := choice.Case("container").(*yang.Container)
			browser.AcceptContainer(visitor, selector, childSchema, t)
		case *yang.List:
			childSchema := choice.Case("list").(*yang.Container)
			browser.AcceptList(visitor, selector, childSchema, t)
		case *yang.Leaf:
			childSchema := choice.Case("leaf").(*yang.Container)
			// TODO: IntLeaf
			visitor.VisitStringLeaf(childSchema, t.(string))
		case *yang.LeafList:
			childSchema := choice.Case("leaf-list").(*yang.Container)
			// TODO: LeafIntList
			visitor.VisitStringLeafList(childSchema, t.([]string))
		case *yang.Choice:
			childSchema := choice.Case("list").(*yang.Container)
			browser.AcceptChoice(visitor, selector, childSchema, t)
		}
		dataDef = dataDef.GetSibling()
	}
}

func (browser *SchemaBrowser) BrowseGrouping(visitor Visitor, selector Selector, schema *yang.Container, grouping yang.DefList) {
	if visitor.EnterList(schema, grouping) {
		field := schema.FirstChild().(yang.Def)
		for field != nil {
			switch field.GetIdent() {
			case :
				default :
					if ! browser.AcceptDefinitionHeader(visitor, selector, field, grouping) {
						if ! browser.AcceptDefList(visitor, selector, field, grouping) {
							// FAIL
						}
					}
			}
		}
		visitor.LeaveContainer(schema, grouping)
	}
}

func (browser *SchemaBrowser) AcceptTypedefs(visitor Visitor, selector Selector, schema *yang.Container, typedefs yang.DefList) {
	if visitor.EnterContainer(schema, typedefs) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			if !browser.AcceptDefList(visitor, selector, field, typedefs) {
				// FAIL
			}
		}
		visitor.LeaveContainer(schema, typedefs)
	}
}

func (browser *SchemaBrowser) AcceptRpcs(visitor Visitor, selector Selector, schema *yang.Container, rpcs yang.DefList) {
	if visitor.EnterContainer(schema, rpcs) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			if !browser.AcceptDefList(visitor, selector, field, rpcs) {
				// FAIL
			}
		}
		visitor.LeaveContainer(schema, rpcs)
	}
}

func (browser *SchemaBrowser) AcceptNotifications(visitor Visitor, selector Selector, schema *yang.Container, notification yang.DefList) {
	if visitor.EnterContainer(schema, notification) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			if !browser.AcceptDefList(visitor, selector, field, notification) {
				// FAIL
			}
		}
		visitor.LeaveContainer(schema, notification)
	}
}

func (browser *SchemaBrowser) Accept(visitor Visitor, selector Selector) {

	m := browser.module
	if visitor.EnterContainer(browser.schema, m) {
		field := browser.schema.GetFirstDef().(yang.Def)
		for field != nil && selector.Selected(field, m) {
			switch field.GetIdent() {
			case "ident", "description":
				browser.AcceptDefinitionHeader(visitor, selector, field, m)
			case "namespace":
				visitor.VisitStringLeaf(field, m.Namespace)
			case "prefix":
				visitor.VisitStringLeaf(field, m.Prefix)

			// LISTS

			case "groupings":
				browser.BrowseGrouping(visitor, selector, field.(*yang.List), m.Groupings)

			case "typedefs":
			// ??
			case "choices":
			// ??
			case "rpcs":
			// ??
			case "notifications":
				// ??
			case "definitions":
				browser.AcceptDataDefinition(visitor, selector, field, browser.module)

			default:
			// FAIL
			}
			field = field.(yang.Containable).GetSibling()
		}
		visitor.LeaveContainer(browser.schema, browser.module)
	}

	browser.module.FirstChild
}

func (browser *SchemaBrowser) AcceptContainer(visitor Visitor, selector Selector, schema *yang.Container, container *yang.Container) {
	if visitor.EnterContainer(schema,container) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			switch schema.GetIdent() {
			case "config":
			case "mandatory":
			case "definitions":
				browser.AcceptDataDefinition(visitor, selector, field, container)
			case "ident", "description":
				browser.AcceptDefinitionHeader(visitor, selector, schema, browser.module)
			default:
				// FAIL
			}
			field = field.(yang.Containable).GetSibling()
		}
		visitor.LeaveContainer(schema, container)
	}
}

func (browser *SchemaBrowser) AcceptList(visitor Visitor, selector Selector, schema *yang.Container, list *yang.List) {
	if visitor.EnterList(schema, list) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			switch schema.GetIdent() {
			case "ident", "description":
				browser.AcceptDefinitionHeader(visitor, selector, schema, list)
			case "config":
			case "mandatory":
			case "choices":
				browser.AcceptChoice(visitor, selector, field, list.GetChoices())
			case "definitions":
				browser.AcceptDataDefinition(visitor, selector, field, list)
			default:
			// FAIL
			}
			field = field.(yang.Containable).GetSibling()
		}
		visitor.LeaveList(schema, list)
	}
}

func (browser *SchemaBrowser) AcceptChoice(visitor Visitor, selector Selector, schema *yang.Container, choice *yang.Choice) {
	if visitor.EnterContainer(schema, choice) {
		field := schema.GetFirstDef().(yang.Def)
		for field != nil {
			switch schema.GetIdent() {
			case "ident", "description":
				browser.AcceptDefinitionHeader(visitor, selector, schema, choice)
			case "choices":
				choiceCase := choice.GetFirstDef()
				for choiceCase != nil {
					browser.AcceptChoiceCase(visitor, selector, field, choiceCase)
					choiceCase = choiceCase.GetSibling()
				}
			default:
			// FAIL
			}
			field = field.(yang.Containable).GetSibling()
		}
		visitor.LeaveChoice(schema, choice)
	}
}

func (browser *SchemaBrowser) AcceptChoiceCase(visitor Visitor, selector Selector, schema *yang.Container, choiceCase *yang.ChoiceCase) {
	if visitor.EnterContainer(schema, choiceCase) {
		field := schema.GetFirstDef()
		for field != nil {
			switch schema.GetIdent() {
			case "ident", "description":
				browser.AcceptDefinitionHeader(visitor, selector, field, choiceCase)
			case "cases":
				browser.AcceptDataDefinition(visitor, selector, field, choiceCase)
				choiceCase := choice.GetFirstDef()
				for choiceCase != nil {
					visitor.R
				}
				browser.AcceptDataDefinition(visitor, selector, field, choice)
			default:
			// FAIL

		}
		visitor.LeaveContainer(schema, choiceCase)
	}
}
