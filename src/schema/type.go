package schema

import (
	"math"
	"strconv"
	"strings"
)


type DataType struct {
	Parent         HasDataType
	Ident          string
	formatPtr      *DataFormat
	rangePtr       *string
	enumeration []string
	minLengthPtr   *int
	maxLengthPtr   *int
	pathPtr        *string
	patternPtr     *string
	defaultPtr     *string
	resolvedPtr    **DataType
	/*
		FractionDigits
		Bit
		Base
		RequireInstance
		Type?!  subtype?
	*/
}

func NewDataType(Parent HasDataType, ident string) (t *DataType) {
	t = &DataType{Parent: Parent, Ident: ident}
	// if not found, then not internal type and Resolve should
	// determine type
	t.SetFormat(DataTypeImplicitFormat(ident))
	return
}

func (y *DataType) resolve() *DataType {
	if y.resolvedPtr == nil {
		var resolved *DataType
		y.resolvedPtr = &resolved
		if y.formatPtr != nil && (*y.formatPtr == FMT_LEAFREF || *y.formatPtr == FMT_LEAFREF_LIST) {
			if y.pathPtr == nil {
				panic("Missing 'path' on leafref " + y.Ident)
			}
			resolvedMeta := FindByPath(y.Parent.GetParent(), *y.pathPtr)
			if resolvedMeta == nil {
				panic("Could not resolve 'path' on leafref " + y.Ident)
			}
			resolved = resolvedMeta.(HasDataType).GetDataType()
		}
		// TODO: else resolve typedefs
	}

	return *y.resolvedPtr
}

func (y *DataType) SetFormat(format DataFormat) {
	if format > 0 {
		y.formatPtr = &format
	}
}

func (y *DataType) Format() (format DataFormat) {
	if y.formatPtr != nil && *y.formatPtr != FMT_LEAFREF && *y.formatPtr != FMT_LEAFREF_LIST {
		format = *y.formatPtr
	} else if resolved := y.resolve(); resolved != nil {
		format = resolved.Format()
	}
	if _, isLeafList := y.Parent.(*LeafList); isLeafList && format < FMT_BINARY_LIST {
		format += 1024
	}
	return
}

func (y *DataType) SetPath(path string) {
	y.pathPtr = &path
}

func (y *DataType) Path() string {
	if y.pathPtr != nil {
		return *y.pathPtr
	}
	if resolved := y.resolve(); resolved != nil {
		return resolved.Path()
	}
	return ""
}

func (y *DataType) SetMinLength(len int) {
	y.minLengthPtr = &len
}

func (y *DataType) MinLength() int {
	if y.minLengthPtr != nil {
		return *y.minLengthPtr
	}
	if resolved := y.resolve(); resolved != nil {
		return resolved.MinLength()
	}
	return 0
}

func (y *DataType) SetMaxLength(len int) {
	y.maxLengthPtr = &len
}

func (y *DataType) MaxLength() int {
	if y.maxLengthPtr != nil {
		return *y.maxLengthPtr
	}
	if resolved := y.resolve(); resolved != nil {
		resolved.MaxLength()
	}
	return math.MaxInt32
}

func (y *DataType) DecodeLength(encoded string) error {
	/* TODO: Support multiple lengths using "|" */
	segments := strings.Split(encoded, "..")
	if len(segments) == 2 {
		if minLength, minErr := strconv.Atoi(segments[0]); minErr != nil {
			return minErr
		} else {
			y.minLengthPtr = &minLength
		}
		if maxLength, maxErr := strconv.Atoi(segments[1]); maxErr != nil {
			return maxErr
		} else {
			y.maxLengthPtr = &maxLength
		}
	} else {
		if maxLength, maxErr := strconv.Atoi(segments[0]); maxErr != nil {
			return maxErr
		} else {
			y.maxLengthPtr = &maxLength
		}
	}
	return nil
}

func (y *DataType) HasDefault() bool {
	if y.defaultPtr != nil {
		return true
	}
	if resolved := y.resolve(); resolved != nil {
		return resolved.HasDefault()
	}
	return false
}

func (y *DataType) SetDefault(def string) {
	y.defaultPtr = &def
}

func (y *DataType) Default() string {
	if y.defaultPtr != nil {
		return *y.defaultPtr
	}
	if resolved := y.resolve(); resolved != nil {
		return resolved.Default()
	}
	return ""
}

func (y *DataType) AddEnumeration(e string) {
	if len(y.enumeration) == 0 {
		y.enumeration = []string{ e }
	} else {
		y.enumeration = append(y.enumeration, e)
	}
}

func (y *DataType) SetEnumeration(en []string) {
	y.enumeration = en
}

func (y *DataType) Enumeration() []string {
	if len(y.enumeration) > 0 {
		return y.enumeration
	}
	if resolved := y.resolve(); resolved != nil {
		return resolved.Enumeration()
	}
	return y.enumeration
}

