package request

type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Teachers []string `json:"teachers"`
	Student  string   `json:"student"`
	Students []string `json:"students"`
}
