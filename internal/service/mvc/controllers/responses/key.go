package responses

type ResourceType string

type Key struct {
	ID   string       `json:"id"`
	Type ResourceType `json:"type"`
}
