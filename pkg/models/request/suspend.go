package request

type SuspendRequest struct {
	Student string `json:"student" binding:"required"`
}
