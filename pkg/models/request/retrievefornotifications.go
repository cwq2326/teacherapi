package request

type ReceieveForNotificationsRequest struct {
	Teacher      string `json:"teacher" binding:"required"`
	Notification string `json:"notification" binding:"required"`
}
