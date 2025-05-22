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
