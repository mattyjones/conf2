package yang

type YangDef interface {
	GetIdent() string
	SetDescription(string)
	GetDescription() string
}

type YangParentable interface {
	AddChild(YangContainerable) error
	GetFirstChild() YangContainerable
	GetLastChild() YangContainerable
}

type YangContainerable interface {
	SetParent(parent YangParentable)
}

type YangDefBase struct {
	Ident string
	Description string
}

func (b *YangDefBase) String() string {
	return b.Ident
}

func (b *ParentableBase) AddChild(child YangContainerable) error {
	//parentable := b.(YangParentable)
	child.SetParent(b)
	b.FirstChild = child
	if b.LastChild == nil {
		b.LastChild = child
	}
	return nil
}

func (b *ParentableBase) GetFirstChild() YangContainerable {
	return b.FirstChild
}
func (b *ParentableBase) GetLastChild() YangContainerable {
	return b.LastChild
}

func (b *ContainableBase) SetParent(parent YangParentable) {
	b.Parent = parent
	b.Sibling = parent.GetFirstChild()
}

type YangRevision struct {
	YangDefBase
}

func (y *YangRevision) GetIdent() (string) {
	return y.Ident
}
func (y *YangRevision) GetDescription() (string) {
	return y.Description
}
func (y *YangRevision) SetDescription(d string) {
	y.Description = d
}

type ParentableBase struct {
	FirstChild YangContainerable
	LastChild YangContainerable
}

type ContainableBase struct {
	Parent YangParentable
	Sibling YangContainerable
}

type YangModule struct {
	YangDefBase
	ParentableBase
	Namespace string
	Revision *YangRevision
	Prefix string
	Rpc ParentableBase
	Notification ParentableBase
	Grouping ParentableBase
	Typedef ParentableBase
}

func (y *YangModule) GetIdent() (string) {
	return y.Ident
}
func (y *YangModule) GetDescription() (string) {
	return y.Description
}
func (y *YangModule) SetDescription(d string) {
	y.Description = d
}

type yangError struct {
	s string
}

func (e yangError) Error() string {
	return e.s
}

func (y *YangModule) AddChild(child YangContainerable) error {
	switch x := child.(type) {
		case *YangRpc:
			return y.Rpc.AddChild(x)
		case *YangNotification:
			return y.Notification.AddChild(x)
		case *YangGrouping:
			return y.Grouping.AddChild(x)
		case *YangTypedef:
			return y.Typedef.AddChild(x)
		default:
			return y.ParentableBase.AddChild(x)
	}
}

type YangContainer struct {
	YangDefBase
	ParentableBase
	ContainableBase
	IsConfig bool
	IsMandatory bool
}
func (y *YangContainer) GetIdent() (string) {
	return y.Ident
}
func (y *YangContainer) GetDescription() (string) {
	return y.Description
}
func (y *YangContainer) SetDescription(d string) {
	y.Description = d
}

func (y *YangContainer) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangContainer) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangList struct {
	YangDefBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
}
func (y *YangList) GetIdent() (string) {
	return y.Ident
}
func (y *YangList) GetDescription() (string) {
	return y.Description
}
func (y *YangList) SetDescription(d string) {
	y.Description = d
}
func (y *YangList) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangList) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangLeaf struct {
	YangDefBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
	Sibling *YangDef
}
func (y *YangLeaf) GetIdent() (string) {
	return y.Ident
}
func (y *YangLeaf) GetDescription() (string) {
	return y.Description
}
func (y *YangLeaf) SetDescription(d string) {
	y.Description = d
}
func (y *YangLeaf) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangLeaf) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangLeafList struct {
	YangDefBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
	Sibling *YangDef
}
func (y *YangLeafList) GetIdent() (string) {
	return y.Ident
}
func (y *YangLeafList) GetDescription() (string) {
	return y.Description
}
func (y *YangLeafList) SetDescription(d string) {
	y.Description = d
}
func (y *YangLeafList) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangLeafList) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangGrouping struct {
	YangDefBase
	ContainableBase
	ParentableBase
	IsConfig bool
	IsMandatory bool
}
func (y *YangGrouping) GetIdent() (string) {
	return y.Ident
}
func (y *YangGrouping) GetDescription() (string) {
	return y.Description
}
func (y *YangGrouping) SetDescription(d string) {
	y.Description = d
}
func (y *YangGrouping) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangGrouping) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangRpcInput struct {
	ContainableBase
	ParentableBase
}
func (y *YangRpcInput) GetIdent() (string) {
	return "input"
}
func (y *YangRpcInput) GetDescription() (string) {
	return ""
}
func (y *YangRpcInput) SetDescription(d string) {
}
func (y *YangRpcInput) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangRpcInput) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangRpcOutput struct {
	ContainableBase
	ParentableBase
}
func (y *YangRpcOutput) GetIdent() (string) {
	return "output"
}
func (y *YangRpcOutput) GetDescription() (string) {
	return ""
}
func (y *YangRpcOutput) SetDescription(d string) {
}
func (y *YangRpcOutput) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangRpcOutput) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangRpc struct {
	YangDefBase
	ContainableBase
	Input *YangRpcInput
	Output *YangRpcOutput
}
func (y *YangRpc) GetIdent() (string) {
	return y.Ident
}
func (y *YangRpc) GetDescription() (string) {
	return y.Description
}
func (y *YangRpc) SetDescription(d string) {
	y.Description = d
}
func (y *YangRpc) AddChild(child YangContainerable) error {
	switch x := child.(type) {
		case *YangRpcInput:
			y.Input = x
		case *YangRpcOutput:
			y.Output = x
		default:
			return yangError{"Error trying to add child of wrong type to rpc"}
	}
	return nil
}
func (y *YangRpc) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}

type YangNotification struct {
	YangDefBase
	ContainableBase
	ParentableBase
}
func (y *YangNotification) GetIdent() (string) {
	return y.Ident
}
func (y *YangNotification) GetDescription() (string) {
	return y.Description
}
func (y *YangNotification) SetDescription(d string) {
	y.Description = d
}
func (y *YangNotification) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangNotification) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}


type YangTypedef struct {
	YangDefBase
	ContainableBase
	ParentableBase
}
func (y *YangTypedef) GetIdent() (string) {
	return y.Ident
}
func (y *YangTypedef) GetDescription() (string) {
	return y.Description
}
func (y *YangTypedef) SetDescription(d string) {
	y.Description = d
}
func (y *YangTypedef) AddChild(child YangContainerable) error {
	return y.ParentableBase.AddChild(child)
}
func (y *YangTypedef) SetParent(parent YangParentable) {
	y.ContainableBase.SetParent(parent)
}
