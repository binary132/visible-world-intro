package models

// Changeset used for auditing
type Changeset struct {
	Table     string        `json:"table"`
	Changes   []*ChangeItem `json:"changes"`
	Operation string        `json:"operation"`
}

type ChangeItem struct {
	Name   string      `json:"name"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

type Changeable interface {
	AddChange(ch ...*Changeset)
}
