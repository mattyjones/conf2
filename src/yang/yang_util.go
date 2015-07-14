package yang

func FindByIdent(l DefList, ident string) Def {
	child := l.GetFirstDef()
	for child != nil {
		if child.GetIdent() == ident {
			return child
		}
	}
	return nil
}

type yangError struct {
	s string
}

func (err *yangError) Error() string {
	return err.s
}