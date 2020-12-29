package api

type ImportDocumentResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Document string `json:"document"`
}
