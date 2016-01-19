package schema

type DataFormat int

// matches list in browse.h
const (
	FMT_BINARY DataFormat = iota + 1
	FMT_BITS
	FMT_BOOLEAN
	FMT_DECIMAL64
	FMT_ENUMERATION
	FMT_IDENTITYREF
	FMT_INSTANCE_IDENTIFIER
	FMT_INT8
	FMT_INT16
	FMT_INT32
	FMT_INT64
	FMT_LEAFREF
	FMT_STRING
	FMT_UINT8
	FMT_UINT16
	FMT_UINT32
	FMT_UINT64
	FMT_UNION
	FMT_ANYDATA
)

const (
	FMT_BINARY_LIST DataFormat = iota + 1025
	FMT_BITS_LIST
	FMT_BOOLEAN_LIST
	FMT_DECIMAL64_LIST
	FMT_ENUMERATION_LIST
	FMT_IDENTITYREF_LIST
	FMT_INSTANCE_IDENTIFIER_LIST
	FMT_INT8_LIST
	FMT_INT16_LIST
	FMT_INT32_LIST
	FMT_INT64_LIST
	FMT_LEAFREF_LIST
	FMT_STRING_LIST
	FMT_UINT8_LIST
	FMT_UINT16_LIST
	FMT_UINT32_LIST
	FMT_UINT64_LIST
	FMT_UNION_LIST
	FMT_ANYDATA_LIST
)

func IsListFormat(f DataFormat) bool {
	return f >= FMT_BINARY_LIST && f <= FMT_UNION_LIST
}

func DataTypeImplicitFormat(typeIdent string) DataFormat {
	return internalTypes[typeIdent]
}

var internalTypes = map[string]DataFormat{
	"binary":              FMT_BINARY,
	"bits":                FMT_BITS,
	"boolean":             FMT_BOOLEAN,
	"decimal64":           FMT_DECIMAL64,
	"enumeration":         FMT_ENUMERATION,
	"identitydef":         FMT_IDENTITYREF,
	"instance-identifier": FMT_INSTANCE_IDENTIFIER,
	"int8":                FMT_INT8,
	"int16":               FMT_INT16,
	"int32":               FMT_INT32,
	"int64":               FMT_INT64,
	"leafref":             FMT_LEAFREF,
	"string":              FMT_STRING,
	"uint8":               FMT_UINT8,
	"uint16":              FMT_UINT16,
	"uint32":              FMT_UINT32,
	"uint64":              FMT_UINT64,
	"union":               FMT_UNION,
	"any":                 FMT_ANYDATA,
}
