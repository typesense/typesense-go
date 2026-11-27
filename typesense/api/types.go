package api

type ImportDocumentResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Document any    `json:"document"` // on success: map[string]interface{}; on error: string
	Id       string `json:"id"`
}

// AnalyticsRuleCreateResponse is one entry of a batch analytics rule create.
// Error is set when the server rejected that rule.
type AnalyticsRuleCreateResponse struct {
	AnalyticsRule
	Error string `json:"error,omitempty"`
}

type StemmingDictionaryWord struct {
	Root string `json:"root"`
	Word string `json:"word"`
}
