package request

// Structure for "/api/suspend" endpoint request body.
type SuspendRequest struct {
	Student string `json:"student" binding:"required"`
}
