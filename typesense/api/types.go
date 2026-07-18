package api

type ImportDocumentResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Document any    `json:"document"` // on success: map[string]interface{}; on error: string
	Id       string `json:"id"`
}

type StemmingDictionaryWord struct {
	Root string `json:"root"`
	Word string `json:"word"`
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
