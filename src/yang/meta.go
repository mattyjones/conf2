package yang

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
type Meta interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
}

// Examples: Module, Container but not Leaf or LeafList
type MetaList interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
	// MetaList
	AddMeta(Meta) error
	GetFirstMeta() Meta
}

type DataDef interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
	// DataDef
	NextDataDef() DataDef
}

type HasGroupings interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
	// MetaList
	AddMeta(Meta) error
	GetFirstMeta() Meta
	// HasGroupings
	GetGroupings() MetaList
}

type HasTypedefs interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
	// MetaList
	AddMeta(Meta) error
	GetFirstMeta() Meta
	// HasTypedefs
	GetTypedefs() MetaList
}

type HasType interface {
	// Identifiable
	GetIdent() string
	// Meta
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
	// HasType
	Type() string
	SetType(dataType string)
}
///////////////////////
// Base structs
///////////////////////

// MetaList implementation helper(s)
type ListBase struct {
	// Parent? - it's normally in MetaBase
	FirstMeta Meta
	LastMeta Meta
}
func (y *ListBase) LinkMeta(impl MetaList, meta Meta) error {
	meta.SetParent(impl)
	if y.LastMeta != nil {
		y.LastMeta.SetSibling(meta)
	}
	y.LastMeta = meta
	if y.FirstMeta == nil {
		y.FirstMeta = meta
	}
	return nil
}

// Meta implementation helpers
type MetaBase struct {
	Parent MetaList
	Sibling Meta
}

// Meta and MetaList combination helpers
type MetaContainer struct {
	Ident string
	MetaBase
	ListBase
}

// Meta
func (y *MetaContainer) GetIdent() (string) {
	return y.Ident
}
// Meta
func (y *MetaContainer) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *MetaContainer) GetParent() MetaList {
	return y.Parent
}
func (y *MetaContainer) GetSibling() Meta {
	return y.Sibling
}
func (y *MetaContainer) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *MetaContainer) AddMeta(meta Meta) error {
	return y.LinkMeta(y, meta)
}
func (y *MetaContainer) GetFirstMeta() Meta {
	return y.FirstMeta
}

type metaError struct {
	Msg string
}

func (e *metaError) Error() string {
	return e.Msg
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
	MetaBase
	Defs MetaContainer
	Rpcs MetaContainer
	Notifications MetaContainer
	Groupings MetaContainer
	Typedefs MetaContainer
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
// Meta
func (y *Module) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Module) GetParent() MetaList {
	return y.Parent
}
func (y *Module) GetSibling() Meta {
	return y.Sibling
}
func (y *Module) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Module) AddMeta(meta Meta) error {
	switch x := meta.(type) {
	case *Rpc:
		y.Rpcs.SetParent(y)
		return y.Rpcs.LinkMeta(y, x)
	case *Notification:
		y.Notifications.SetParent(y)
		return y.Notifications.LinkMeta(y, x)
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, x)
	case *Typedef:
		y.Typedefs.SetParent(y)
		return y.Typedefs.LinkMeta(y, x)
	default:
		y.Defs.SetParent(y)
		return y.Defs.LinkMeta(y, x)
	}
}
// technically not true, it's the MetaContainers, but we'll see how this pans out
func (y *Module) GetFirstMeta() Meta {
	return y.Defs.GetFirstMeta()
}
func (y *Module) DataDefs() MetaList {
	return &y.Defs
}
func (y *Module) GetRpcs() MetaList {
	return &y.Rpcs
}
func (y *Module) GetNotifications() MetaList {
	return &y.Notifications
}
// HasGroupings
func (y *Module) GetGroupings() MetaList {
	return &y.Groupings
}
func (y *Module) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type ChoiceDecider func(Choice, ChoiceCase, interface{})


type Choice struct {
	Ident string
	Description string
	MetaBase
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
// Meta
func (y *Choice) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Choice) GetParent() MetaList {
	return y.Parent
}
func (y *Choice) GetSibling() Meta {
	return y.Sibling
}
func (y *Choice) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Choice) AddMeta(meta Meta) error {
	return y.LinkMeta(y, meta)
}
func (y *Choice) GetFirstMeta() Meta {
	return y.FirstMeta
}
// Other
func (c *Choice) GetCase(ident string) *ChoiceCase {
	return FindByPathWithoutResolvingProxies(c, ident).(*ChoiceCase)
}
// MetaProxy
func (y *Choice) ResolveProxy() MetaIterator {
	return &MetaListIterator{position:y.GetFirstMeta(),resolveProxies:true}
}

////////////////////////////////////////////////////

type ChoiceCase struct {
	Ident string
	MetaBase
	ListBase
}
// Identifiable
func (y *ChoiceCase) GetIdent() (string) {
	return y.Ident
}
// Meta
func (y *ChoiceCase) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *ChoiceCase) GetParent() MetaList {
	return y.Parent
}
func (y *ChoiceCase) GetSibling() Meta {
	return y.Sibling
}
func (y *ChoiceCase) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *ChoiceCase) AddMeta(meta Meta) error {
	return y.LinkMeta(y, meta)
}
func (y *ChoiceCase) GetFirstMeta() Meta {
	return y.FirstMeta
}
// MetaProxy
func (y *ChoiceCase) ResolveProxy() MetaIterator {
	return &MetaListIterator{position:y.GetFirstMeta(), resolveProxies:true}
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
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs MetaContainer
	Config bool
	Mandatory bool
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
// Meta
func (y *Container) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Container) GetParent() MetaList {
	return y.Parent
}
func (y *Container) GetSibling() Meta {
	return y.Sibling
}
func (y *Container) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Container) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, meta)
	default:
		e := y.LinkMeta(y, meta)
		return e
	}
}
func (y *Container) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *Container) GetGroupings() MetaList {
	return &y.Groupings
}
// HasTypedefs
func (y *Container) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type List struct {
	Ident string
	Description string
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs MetaContainer
	Config bool
	Mandatory bool
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
// Meta
func (y *List) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *List) GetParent() MetaList {
	return y.Parent
}
func (y *List) GetSibling() Meta {
	return y.Sibling
}
func (y *List) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *List) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, meta)
	default:
		return y.LinkMeta(y, meta)
	}
}
func (y *List) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *List) GetGroupings() MetaList {
	return &y.Groupings
}
// HasTypedefs
func (y *List) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Leaf struct {
	Ident string
	Description string
	MetaBase
	Config bool
	Mandatory bool
	DataType string
}

// Distinguishes the concrete type in choice-cases
func (y *Leaf) Leaf() Meta {
	return y
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
// Meta
func (y *Leaf) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Leaf) GetParent() MetaList {
	return y.Parent
}
func (y *Leaf) GetSibling() Meta {
	return y.Sibling
}
func (y *Leaf) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// HasType
func (y *Leaf) Type() string {
	return y.DataType
}
func (y *Leaf) SetType(dataType string) {
	y.DataType = dataType
}


////////////////////////////////////////////////////

type LeafList struct {
	Ident string
	Description string
	MetaBase
	Config bool
	Mandatory bool
	DataType string
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
// Meta
func (y *LeafList) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *LeafList) GetParent() MetaList {
	return y.Parent
}
func (y *LeafList) GetSibling() Meta {
	return y.Sibling
}
func (y *LeafList) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// HasType
func (y *LeafList) Type() string {
	return y.DataType
}
func (y *LeafList) SetType(dataType string) {
	y.DataType = dataType
}

////////////////////////////////////////////////////

type Grouping struct {
	Ident string
	Description string
	MetaBase
	ListBase
	Config bool
	Mandatory bool
	Groupings MetaContainer
	Typedefs MetaContainer
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
// Meta
func (y *Grouping) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Grouping) GetParent() MetaList {
	return y.Parent
}
func (y *Grouping) GetSibling() Meta {
	return y.Sibling
}
func (y *Grouping) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Grouping) AddMeta(meta Meta) error {
	return y.LinkMeta(y, meta)
}
func (y *Grouping) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *Grouping) GetGroupings() MetaList {
	return &y.Groupings
}
// HasTypedefs
func (y *Grouping) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type RpcInput struct {
	MetaBase
	ListBase
	Typedefs MetaContainer
	Groupings MetaContainer
}
// Identifiable
func (y *RpcInput) GetIdent() (string) {
	// Not technically true, but works
	return "input"
}
// Meta
func (y *RpcInput) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *RpcInput) GetParent() MetaList {
	return y.Parent
}
func (y *RpcInput) GetSibling() Meta {
	return y.Sibling
}
func (y *RpcInput) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *RpcInput) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, meta)
	default:
		return y.LinkMeta(y, meta)
	}
}
func (y *RpcInput) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *RpcInput) GetGroupings() MetaList {
	return &y.Groupings
}
// HasTypedefs
func (y *RpcInput) GetTypedefs() MetaList {
	return &y.Typedefs
}


////////////////////////////////////////////////////

type RpcOutput struct {
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs MetaContainer
}
// Identifiable
func (y *RpcOutput) GetIdent() (string) {
	return "output"
}
// Meta
func (y *RpcOutput) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *RpcOutput) GetParent() MetaList {
	return y.Parent
}
func (y *RpcOutput) GetSibling() Meta {
	return y.Sibling
}
func (y *RpcOutput) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *RpcOutput) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, meta)
	default:
		return y.LinkMeta(y, meta)
	}
}
func (y *RpcOutput) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *RpcOutput) GetGroupings() MetaList {
	return &y.Groupings
}
// HasTypedefs
func (y *RpcOutput) GetTypedefsGroupings() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Rpc struct {
	Ident string
	Description string
	MetaBase
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
// Meta
func (y *Rpc) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Rpc) GetParent() MetaList {
	return y.Parent
}
func (y *Rpc) GetSibling() Meta {
	return y.Sibling
}
func (y *Rpc) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Rpc) AddMeta(meta Meta) error {
	switch t:= meta.(type) {
	case *RpcInput:
		t.SetParent(y)
		y.Input = t
		return nil
	case *RpcOutput:
		t.SetParent(y)
		y.Output = t
	default:
		return &metaError{"Illegal call to add metainition: rpc has fixed input and output children"}
	}
	if y.Output != nil {
		y.Input.Sibling = y.Output
	}
	return nil
}
func (y *Rpc) GetFirstMeta() Meta {
	if y.Input != nil {
		return y.Input
	}
	return y.Output
}

////////////////////////////////////////////////////

type Notification struct {
	Ident     string
	Description string
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs MetaContainer
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
// Meta
func (y *Notification) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Notification) GetParent() MetaList {
	return y.Parent
}
func (y *Notification) GetSibling() Meta {
	return y.Sibling
}
func (y *Notification) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Notification) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.LinkMeta(y, meta)
	default:
		return y.LinkMeta(y, meta)
	}
}
func (y *Notification) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasGroupings
func (y *Notification) GetGroupings() MetaList {
	return &y.Groupings
}
// HasGroupings
func (y *Notification) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Typedef struct {
	Ident string
	Description string
	MetaContainer
	DataType string
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
// Meta
func (y *Typedef) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Typedef) GetParent() MetaList {
	return y.Parent
}
func (y *Typedef) GetSibling() Meta {
	return y.Sibling
}
func (y *Typedef) SetSibling(sibling Meta) {
	y.Sibling = sibling
}
// MetaList
func (y *Typedef) AddMeta(meta Meta) error {
	y.MetaContainer.SetParent(y)
	return y.MetaContainer.LinkMeta(y, meta)
}
func (y *Typedef) GetFirstMeta() Meta {
	return y.FirstMeta
}
// HasType
func (y *Typedef) Type() string {
	return y.DataType
}
func (y *Typedef) SetType(dataType string) {
	y.DataType = dataType
}
// Typedef
func (y *Typedef) GetEnumerations() MetaList {
	// TODO
	return nil
}

////////////////////////////////////////////////////

type Uses struct {
	Ident string
	Description string
	MetaBase
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
// Meta
func (y *Uses) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Uses) GetParent() MetaList {
	return y.Parent
}
func (y *Uses) GetSibling() Meta {
	return y.Sibling
}
func (y *Uses) SetSibling(sibling Meta) {
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
// MetaProxy
func (y *Uses) ResolveProxy() MetaIterator {
	if g := y.FindGrouping(y.Ident); g != nil {
		return NewMetaListIterator(g, true)
	}
	return nil
}
