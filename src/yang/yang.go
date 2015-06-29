package yang

type Definition interface {
	GetIdent() string
	SetDescription(string)
	GetDescription() string
}

type Parentable interface {
	AddChild(Containable) error
	GetFirstChild() Containable
	GetLastChild() Containable
}

type Containable interface {
	SetParent(parent Parentable)
}

type DefinitionBase struct {
	Ident string
	Description string
}

func (b *DefinitionBase) String() string {
	return b.Ident
}

func (b *ParentableBase) AddChild(child Containable) error {
	//parentable := b.(Parentable)
	child.SetParent(b)
	b.FirstChild = child
	if b.LastChild == nil {
		b.LastChild = child
	}
	return nil
}

func (b *ParentableBase) GetFirstChild() Containable {
	return b.FirstChild
}
func (b *ParentableBase) GetLastChild() Containable {
	return b.LastChild
}

func (b *ContainableBase) SetParent(parent Parentable) {
	b.Parent = parent
	b.Sibling = parent.GetFirstChild()
}

type Revision struct {
	DefinitionBase
}

func (y *Revision) GetIdent() (string) {
	return y.Ident
}
func (y *Revision) GetDescription() (string) {
	return y.Description
}
func (y *Revision) SetDescription(d string) {
	y.Description = d
}

type ParentableBase struct {
	FirstChild Containable
	LastChild Containable
}

type ContainableBase struct {
	Parent Parentable
	Sibling Containable
}

type Module struct {
	DefinitionBase
	ParentableBase
	Namespace string
	Revision *Revision
	Prefix string
	Rpc ParentableBase
	Notification ParentableBase
	Grouping ParentableBase
	Typedef ParentableBase
}

func (y *Module) GetIdent() (string) {
	return y.Ident
}
func (y *Module) GetDescription() (string) {
	return y.Description
}
func (y *Module) SetDescription(d string) {
	y.Description = d
}

type yangError struct {
	s string
}

func (e yangError) Error() string {
	return e.s
}

func (y *Module) AddChild(child Containable) error {
	switch x := child.(type) {
		case *Rpc:
			return y.Rpc.AddChild(x)
		case *Notification:
			return y.Notification.AddChild(x)
		case *Grouping:
			return y.Grouping.AddChild(x)
		case *Typedef:
			return y.Typedef.AddChild(x)
		default:
			return y.ParentableBase.AddChild(x)
	}
}

type Container struct {
	DefinitionBase
	ParentableBase
	ContainableBase
	IsConfig bool
	IsMandatory bool
}
func (y *Container) GetIdent() (string) {
	return y.Ident
}
func (y *Container) GetDescription() (string) {
	return y.Description
}
func (y *Container) SetDescription(d string) {
	y.Description = d
}

func (y *Container) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *Container) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type List struct {
	DefinitionBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
}
func (y *List) GetIdent() (string) {
	return y.Ident
}
func (y *List) GetDescription() (string) {
	return y.Description
}
func (y *List) SetDescription(d string) {
	y.Description = d
}
func (y *List) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *List) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type Leaf struct {
	DefinitionBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
	Sibling *Definition
}
func (y *Leaf) GetIdent() (string) {
	return y.Ident
}
func (y *Leaf) GetDescription() (string) {
	return y.Description
}
func (y *Leaf) SetDescription(d string) {
	y.Description = d
}
func (y *Leaf) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *Leaf) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type LeafList struct {
	DefinitionBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
	Sibling *Definition
}
func (y *LeafList) GetIdent() (string) {
	return y.Ident
}
func (y *LeafList) GetDescription() (string) {
	return y.Description
}
func (y *LeafList) SetDescription(d string) {
	y.Description = d
}
func (y *LeafList) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *LeafList) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type Grouping struct {
	DefinitionBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
}
func (y *Grouping) GetIdent() (string) {
	return y.Ident
}
func (y *Grouping) GetDescription() (string) {
	return y.Description
}
func (y *Grouping) SetDescription(d string) {
	y.Description = d
}
func (y *Grouping) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *Grouping) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type RpcInput struct {
	ContainableBase
	ParentableBase
}
func (y *RpcInput) GetIdent() (string) {
	return "input"
}
func (y *RpcInput) GetDescription() (string) {
	return ""
}
func (y *RpcInput) SetDescription(d string) {
}
func (y *RpcInput) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *RpcInput) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type RpcOutput struct {
	ContainableBase
	ParentableBase
}
func (y *RpcOutput) GetIdent() (string) {
	return "output"
}
func (y *RpcOutput) GetDescription() (string) {
	return ""
}
func (y *RpcOutput) SetDescription(d string) {
}
func (y *RpcOutput) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *RpcOutput) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type Rpc struct {
	DefinitionBase
	ContainableBase
	Input *RpcInput
	Output *RpcOutput
}
func (y *Rpc) GetIdent() (string) {
	return y.Ident
}
func (y *Rpc) GetDescription() (string) {
	return y.Description
}
func (y *Rpc) SetDescription(d string) {
	y.Description = d
}
func (y *Rpc) AddChild(child Containable) error {
	switch x := child.(type) {
		case *RpcInput:
			y.Input = x
		case *RpcOutput:
			y.Output = x
		default:
			return yangError{"Error trying to add child of wrong type to rpc"}
	}
	return nil
}
func (y *Rpc) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}

type Notification struct {
	DefinitionBase
	ContainableBase
	ParentableBase
}
func (y *Notification) GetIdent() (string) {
	return y.Ident
}
func (y *Notification) GetDescription() (string) {
	return y.Description
}
func (y *Notification) SetDescription(d string) {
	y.Description = d
}
func (y *Notification) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *Notification) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}


type Typedef struct {
	DefinitionBase
	ContainableBase
	ParentableBase
}
func (y *Typedef) GetIdent() (string) {
	return y.Ident
}
func (y *Typedef) GetDescription() (string) {
	return y.Description
}
func (y *Typedef) SetDescription(d string) {
	y.Description = d
}
func (y *Typedef) AddChild(child Containable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *Typedef) SetParent(parent Parentable) {
	y.ContainableBase.SetParent(parent)
}
