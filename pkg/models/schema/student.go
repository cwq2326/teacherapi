package schema

// Schema for student relation.
type Student struct {
	Email     string `json:"email"`
	Suspended bool   `json:"student"`
}
