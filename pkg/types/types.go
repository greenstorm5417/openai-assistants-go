package types

type Metadata map[string]interface{}

type FileRef struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}