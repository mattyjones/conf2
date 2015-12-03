package schema

// Tri-state boolean - true, false, and undeclared
type Flag int

const (
	UNSET Flag = iota
	SET_TRUE
	SET_FALSE
)

func (f Flag) IsSet() bool {
	return f != UNSET
}

func (f Flag) Bool() bool {
	return f == SET_TRUE
}

func (f Flag) Set(b bool) {
	if b {
		f = SET_TRUE
	} else {
		f = SET_FALSE
	}
}

type Details struct {
	ConfigFlag    Flag
	MandatoryFlag Flag
}

func (d *Details) Config(p *Path) bool {
	switch d.ConfigFlag {
	case SET_TRUE:
		return true
	case SET_FALSE:
		return false
	}

	// if details are on leaf, then p is parent container, otherwise
	// p is what we're supposed to check config for
	if p != nil {
		if hasDetails, ok := p.meta.(HasDetails); ok {
			if parentDetails := hasDetails.Details(); parentDetails == d {
				return hasDetails.Details().Config(p.parent)
			}
			return hasDetails.Details().Config(p)
		}
	}
	return true
}

func (d *Details) SetConfig(config bool) {
	if config {
		d.ConfigFlag = SET_TRUE
	} else {
		d.ConfigFlag = SET_FALSE
	}
}

func (d *Details) Mandatory() bool {
	return SET_TRUE == d.MandatoryFlag
}
