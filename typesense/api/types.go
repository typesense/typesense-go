package api

type ImportDocumentResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Document string `json:"document"`
	Id       string `json:"id"`
}

type StemmingDictionaryWord struct {
	Root string `json:"root"`
	Word string `json:"word"`
}
