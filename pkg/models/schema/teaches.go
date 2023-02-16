package schema

// Schema for teaches relation.
type Teaches struct {
	Teacher string `json:"teacher"`
	Student string `json:"student"`
}
