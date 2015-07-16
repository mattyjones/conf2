package yang

import "log"

///////////////////
// Interfaces
//////////////////

// Examples: Just about everything
type Identifiable interface {
	GetIdent() string
}

// Examples: Most everything except items like ChoiceCase, RpcInput
type Describable interface {
	GetIdent() string
	SetDescription(string)
	GetDescription() string
}

// Examples: Things that have more than one.
type Def interface {
	// Identifiable
	GetIdent() string
	// Def
	GetParent() DefList
	SetParent(DefList)
	GetSibling() Def
	SetSibling(Def)
}

// Examples: Module, Container but not Leaf or LeafList
type DefList interface {
	// Def
	GetParent() DefList
	SetParent(DefList)
	GetSibling() Def
	SetSibling(Def)
	// DefList
	AddDef(Def) error
	GetFirstDef() Def
}

type DataDef interface {
	// Identifiable
	GetIdent() string
	// Def
	GetParent() DefList
	SetParent(DefList)
	GetSibling() Def
	SetSibling(Def)
	// DataDef
	NextDataDef() DataDef
}

// Examples: Container, List, Grouping
type HasChoices interface {
	// Identifiable
	GetIdent() string
	// Def
	GetParent() DefList
	SetParent(DefList)
	GetSibling() Def
	SetSibling(Def)
	// DefList
	AddDef(Def) error
	GetFirstDef() Def
	// HasChoices
	GetFirstChoice() *Choice
}

type HasGroupings interface {
	// Identifiable
	GetIdent() string
	// Def
	GetParent() DefList
	SetParent(DefList)
	GetSibling() Def
	SetSibling(Def)
	// DefList
	AddDef(Def) error
	GetFirstDef() Def
	// HasGroupings
	GetGroupings() DefList
}

///////////////////////
// Base structs
///////////////////////

// DefList implementation helper(s)
type ListBase struct {
	// Parent? - it's normally in DefBase
	FirstDef Def
	LastDef Def
}
func (y *ListBase) LinkDef(impl DefList, def Def) error {
	def.SetParent(impl)
	if y.LastDef != nil {
		y.LastDef.SetSibling(def)
	}
	y.LastDef = def
	if y.FirstDef == nil {
		y.FirstDef = def
	}
	return nil
}

// Def implementation helpers
type DefBase struct {
	Parent DefList
	Sibling Def
}
//func (y *DefBase) SetParent(parent DefList) {
//	y.Parent = parent
	//y.Sibling = parent.GetFirstDef()
//}

// Def and DefList combination helpers
type DefContainer struct {
	DefBase
	ListBase
}
//func (y *DefContainer) LinkDef(impl DefList, def Def) error {
//	return y.ListBase.LinkDef(impl, def)
//}

// Store DefList
func (y *DefContainer) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *DefContainer) GetParent() DefList {
	return y.DefBase.Parent
}
func (y *DefContainer) AddDef(def Def) error {
	return y.LinkDef(y, def)
}
func (y *DefContainer) GetFirstDef() Def {
	return y.FirstDef
}
func (y *DefContainer) GetSibling() Def {
	return y.FirstDef
}
func (y *DefContainer) SetSibling(sibling Def) {
	y.DefBase.Sibling = sibling
}

////////////////////////
// Implementations
/////////////////////////

type Module struct {
	Ident string
	Description string
	Namespace string
	Revision *Revision
	Prefix string
	DefBase
	Defs DefContainer
	Rpcs DefContainer
	Notifications DefContainer
	Groupings DefContainer
	Choices DefContainer
	Typedefs DefContainer
}
// Identifiable
func (y *Module) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Module) GetDescription() (string) {
	return y.Description
}
func (y *Module) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Module) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Module) GetParent() DefList {
	return y.Parent
}
func (y *Module) GetSibling() Def {
	return y.Sibling
}
func (y *Module) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Module) AddDef(def Def) error {
	switch x := def.(type) {
	case *Rpc:
		y.Rpcs.SetParent(y)
		return y.Rpcs.LinkDef(y, x)
	case *Notification:
		y.Notifications.SetParent(y)
		return y.Notifications.LinkDef(y, x)
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, x)
	case *Typedef:
		y.Typedefs.SetParent(y)
		return y.Typedefs.LinkDef(y, x)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, x)
	default:
		y.Defs.SetParent(y)
		return y.Defs.LinkDef(y, x)
	}
}
// technically not true, it's the DefContainers, but we'll see how this pans out
func (y *Module) GetFirstDef() Def {
	return y.Defs.GetFirstDef()
}
func (y *Module) DataDefs() DefList {
	return &y.Defs
}
func (y *Module) GetRpcs() DefList {
	return &y.Rpcs
}
func (y *Module) GetNotifications() DefList {
	return &y.Notifications
}
// HasGroupings
func (y *Module) GetGroupings() DefList {
	log.Println("getting module groupings", y.Groupings)
	return &y.Groupings
}
func (y *Module) GetChoices() DefList {
	return &y.Choices
}
func (y *Module) GetTypeDefs() DefList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Choice struct {
	Ident string
	Description string
	DefBase
	ListBase
}
// Identifiable
func (y *Choice) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Choice) GetDescription() (string) {
	return y.Description
}
func (y *Choice) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Choice) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Choice) GetParent() DefList {
	return y.Parent
}
func (y *Choice) GetSibling() Def {
	return y.Sibling
}
func (y *Choice) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Choice) AddDef(def Def) error {
	return y.LinkDef(y, def)
}
func (y *Choice) GetFirstDef() Def {
	return y.FirstDef
}
// Other
func (c *Choice) GetCase(ident string) *ChoiceCase {
	return FindByPathWithoutResolvingProxies(c, ident).(*ChoiceCase)
}
// DefProxy
func (y *Choice) ResolveProxy() DefIterator {
	return &DefListIterator{position:y.GetFirstDef(),resolveProxies:true}
}

////////////////////////////////////////////////////

type ChoiceCase struct {
	Ident string
	DefBase
	ListBase
}
// Identifiable
func (y *ChoiceCase) GetIdent() (string) {
	return y.Ident
}
// Def
func (y *ChoiceCase) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *ChoiceCase) GetParent() DefList {
	return y.Parent
}
func (y *ChoiceCase) GetSibling() Def {
	return y.Sibling
}
func (y *ChoiceCase) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *ChoiceCase) AddDef(def Def) error {
	return y.LinkDef(y, def)
}
func (y *ChoiceCase) GetFirstDef() Def {
	return y.FirstDef
}
// DefProxy
func (y *ChoiceCase) ResolveProxy() DefIterator {
	return &DefListIterator{position:y.GetFirstDef(), resolveProxies:true}
}

////////////////////////////////////////////////////

type Revision struct {
	Ident string
	Description string
}
// Identifiable
func (y *Revision) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Revision) GetDescription() (string) {
	return y.Description
}
func (y *Revision) SetDescription(d string) {
	y.Description = d
}

////////////////////////////////////////////////////

type Container struct {
	Ident string
	Description string
	DefBase
	ListBase
	Groupings DefContainer
	Choices DefContainer
	IsConfig bool
	IsMandatory bool
}
// Identifiable
func (y *Container) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Container) GetDescription() (string) {
	return y.Description
}
func (y *Container) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Container) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Container) GetParent() DefList {
	return y.Parent
}
func (y *Container) GetSibling() Def {
	return y.Sibling
}
func (y *Container) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Container) AddDef(def Def) error {
	switch def.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, def)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		e := y.LinkDef(y, def)
		return e
	}
}
func (y *Container) GetFirstDef() Def {
	return y.FirstDef
}
// HasChoices
func (y *Container) GetFirstChoice() *Choice {
	y.Choices.SetParent(y)
	return y.Choices.FirstDef.(*Choice)
}
// HasGroupings
func (y *Container) GetGroupings() DefList {
	return &y.Groupings
}

////////////////////////////////////////////////////

type List struct {
	Ident string
	Description string
	DefBase
	ListBase
	Groupings DefContainer
	Choices DefContainer
	IsConfig bool
	IsMandatory bool
}
// Identifiable
func (y *List) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *List) GetDescription() (string) {
	return y.Description
}
func (y *List) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *List) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *List) GetParent() DefList {
	return y.Parent
}
func (y *List) GetSibling() Def {
	return y.Sibling
}
func (y *List) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *List) AddDef(def Def) error {
	switch def.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, def)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		return y.LinkDef(y, def)
	}
}
func (y *List) GetFirstDef() Def {
	return y.FirstDef
}
// HasChoices
func (y *List) GetFirstChoice() *Choice {
	y.Choices.SetParent(y)
	return y.Choices.FirstDef.(*Choice)
}
// HasGroupings
func (y *List) GetGroupings() DefList {
	return &y.Groupings
}

////////////////////////////////////////////////////

type Leaf struct {
	Ident string
	Description string
	DefBase
	IsConfig bool
	IsMandatory bool
}
// Identifiable
func (y *Leaf) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Leaf) GetDescription() (string) {
	return y.Description
}
func (y *Leaf) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Leaf) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Leaf) GetParent() DefList {
	return y.Parent
}
func (y *Leaf) GetSibling() Def {
	return y.Sibling
}
func (y *Leaf) SetSibling(sibling Def) {
	y.Sibling = sibling
}

////////////////////////////////////////////////////

type LeafList struct {
	Ident string
	Description string
	DefBase
	IsConfig bool
	IsMandatory bool
}
// Identifiable
func (y *LeafList) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *LeafList) GetDescription() (string) {
	return y.Description
}
func (y *LeafList) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *LeafList) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *LeafList) GetParent() DefList {
	return y.Parent
}
func (y *LeafList) GetSibling() Def {
	return y.Sibling
}
func (y *LeafList) SetSibling(sibling Def) {
	y.Sibling = sibling
}

////////////////////////////////////////////////////

type Grouping struct {
	Ident string
	Description string
	DefBase
	ListBase
	Choices DefContainer
	IsConfig bool
	IsMandatory bool
}
// Identifiable
func (y *Grouping) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Grouping) GetDescription() (string) {
	return y.Description
}
func (y *Grouping) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Grouping) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Grouping) GetParent() DefList {
	return y.Parent
}
func (y *Grouping) GetSibling() Def {
	return y.Sibling
}
func (y *Grouping) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Grouping) AddDef(def Def) error {
	switch def.(type) {
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		return y.LinkDef(y, def)
	}
}
func (y *Grouping) GetFirstDef() Def {
	return y.FirstDef
}
// HasChoices
func (y *Grouping) GetFirstChoice() *Choice {
	y.Choices.SetParent(y)
	return y.Choices.FirstDef.(*Choice)
}
// DefProxy
//func (y *Grouping) ResolveProxy() DefIterator {
//	return &DefListIterator{position:y.GetFirstDef(), resolveProxies:true}
//}

////////////////////////////////////////////////////

type RpcInput struct {
	DefBase
	ListBase
	Groupings DefContainer
	Choices DefContainer
}
// Identifiable
func (y *RpcInput) GetIdent() (string) {
	// Not technically true, but works
	return "input"
}
// Def
func (y *RpcInput) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *RpcInput) GetParent() DefList {
	return y.Parent
}
func (y *RpcInput) GetSibling() Def {
	return y.Sibling
}
func (y *RpcInput) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *RpcInput) AddDef(def Def) error {
	switch def.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, def)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		return y.LinkDef(y, def)
	}
}
func (y *RpcInput) GetFirstDef() Def {
	return y.FirstDef
}
// HasChoices
func (y *RpcInput) GetFirstChoice() *Choice {
	y.Choices.SetParent(y)
	return y.Choices.FirstDef.(*Choice)
}
// HasGroupings
func (y *RpcInput) GetGroupings() DefList {
	return &y.Groupings
}


////////////////////////////////////////////////////

type RpcOutput struct {
	DefBase
	ListBase
	Groupings DefContainer
	Choices DefContainer
}
// Identifiable
func (y *RpcOutput) GetIdent() (string) {
	return "output"
}
// Def
func (y *RpcOutput) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *RpcOutput) GetParent() DefList {
	return y.Parent
}
func (y *RpcOutput) GetSibling() Def {
	return y.Sibling
}
func (y *RpcOutput) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *RpcOutput) AddDef(def Def) error {
	switch def.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, def)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		return y.LinkDef(y, def)
	}
}
func (y *RpcOutput) GetFirstDef() Def {
	return y.FirstDef
}
// HasChoices
func (y *RpcOutput) GetFirstChoice() *Choice {
	y.Choices.SetParent(y)
	return y.Choices.FirstDef.(*Choice)
}
// HasGroupings
func (y *RpcOutput) GetGroupings() DefList {
	return &y.Groupings
}

////////////////////////////////////////////////////

type Rpc struct {
	Ident string
	Description string
	DefBase
	Input *RpcInput
	Output *RpcOutput
}
// Identifiable
func (y *Rpc) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Rpc) GetDescription() (string) {
	return y.Description
}
func (y *Rpc) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Rpc) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Rpc) GetParent() DefList {
	return y.Parent
}
func (y *Rpc) GetSibling() Def {
	return y.Sibling
}
func (y *Rpc) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Rpc) AddDef(def Def) error {
	switch t:= def.(type) {
	case *RpcInput:
		t.SetParent(y)
		y.Input = t
		return nil
	case *RpcOutput:
		t.SetParent(y)
		y.Output = t
	default:
		return &yangError{"Illegal call to add definition: rpc has fixed input and output children"}
	}
	if y.Output != nil {
		y.Input.Sibling = y.Output
	}
	return nil
}
func (y *Rpc) GetFirstDef() Def {
	if y.Input != nil {
		return y.Input
	}
	return y.Output
}

////////////////////////////////////////////////////

type Notification struct {
	Ident string
	Description string
	DefBase
	ListBase
	Groupings DefContainer
	Choices DefContainer
}
// Identifiable
func (y *Notification) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Notification) GetDescription() (string) {
	return y.Description
}
func (y *Notification) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Notification) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Notification) GetParent() DefList {
	return y.Parent
}
func (y *Notification) GetSibling() Def {
	return y.Sibling
}
func (y *Notification) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Notification) AddDef(def Def) error {
	switch def.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkDef(y, def)
	case *Choice:
		y.Choices.SetParent(y)
		return y.Choices.LinkDef(y, def)
	default:
		return y.LinkDef(y, def)
	}
}
func (y *Notification) GetFirstDef() Def {
	return y.FirstDef
}

////////////////////////////////////////////////////

type Typedef struct {
	Ident string
	Description string
	DefContainer
}
// Identifiable
func (y *Typedef) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Typedef) GetDescription() (string) {
	return y.Description
}
func (y *Typedef) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Typedef) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Typedef) GetParent() DefList {
	return y.Parent
}
func (y *Typedef) GetSibling() Def {
	return y.Sibling
}
func (y *Typedef) SetSibling(sibling Def) {
	y.Sibling = sibling
}
// DefList
func (y *Typedef) AddDef(def Def) error {
	y.DefContainer.SetParent(y)
	return y.DefContainer.LinkDef(y, def)
}
func (y *Typedef) GetFirstDef() Def {
	return y.FirstDef
}

////////////////////////////////////////////////////

type Uses struct {
	Ident string
	Description string
	DefBase
	grouping *Grouping
	// augment
	// if-feature
	// refine
	// reference
	// status
	// when
}
// Identifiable
func (y *Uses) GetIdent() (string) {
	return y.Ident
}
// Describable
func (y *Uses) GetDescription() (string) {
	return y.Description
}
func (y *Uses) SetDescription(d string) {
	y.Description = d
}
// Def
func (y *Uses) SetParent(parent DefList) {
	y.Parent = parent
}
func (y *Uses) GetParent() DefList {
	return y.Parent
}
func (y *Uses) GetSibling() Def {
	return y.Sibling
}
func (y *Uses) SetSibling(sibling Def) {
	y.Sibling = sibling
}
func (y *Uses) FindGrouping(ident string) *Grouping {
	// lazy load grouping
	if y.grouping == nil {
		p := y.GetParent()
		for p != nil && y.grouping == nil {
			if withGrouping, hasGrouping := p.(HasGroupings); hasGrouping {
				found := FindByPath(withGrouping.GetGroupings(), y.GetIdent())
				if found != nil {
					y.grouping = found.(*Grouping)
				}
			}
			p  = p.GetParent()
		}
	}
	return y.grouping
}
// DefProxy
func (y *Uses) ResolveProxy() DefIterator {
	if g := y.FindGrouping(y.Ident); g != nil {
		return NewDefListIterator(g, true)
	}
	return nil
}
