package schema

// Tri-state boolean - true, false, and not-defined
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

func (d *Details) Config(p *MetaPath) bool {
	switch d.ConfigFlag {
	case SET_TRUE:
		return true
	case SET_FALSE:
		return false
	}
	if p.ParentPath != nil {
		if hasDetails, ok := p.Parent().(HasDetails); ok {
			return hasDetails.Details().Config(p.ParentPath)
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
