package request

// Structure for "/api/register" endpoint request body.
type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Teachers []string `json:"teachers"`
	Student  string   `json:"student"`
	Students []string `json:"students"`
}
