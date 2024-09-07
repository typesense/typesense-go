package api

type ImportDocumentResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Document string `json:"document"`
}

type Type string

// Enum for all Types in Typesense
const (
	STRING        Type = "string"
	STRINGARRAY        = "string[]"
	INT32              = "int32"
	INT32ARRAY         = "int32[]"
	INT64              = "int64"
	INT64ARRAY         = "int64[]"
	FLOAT              = "float"
	FLOATARRAY         = "float[]"
	BOOL               = "bool"
	BOOLARRAY          = "bool[]"
	GEOPOINT           = "geopoint"
	GEOPOINTARRAY      = "geopoint[]"
	OBJECT             = "object" // object is comparable to a go struct
	OBJECTARRAY        = "object[]"
	STRINGPTR          = "string*" // special type that can be string or []string
)
