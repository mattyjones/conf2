package comm

// Before investing too much time into data serialization, consider using flatbuffers
// library to r/w data structures.

const CSTR_TERM = byte(0)

const (
	TRUE_BYTE  = byte(1)
	FALSE_BYTE = byte(0)
)
