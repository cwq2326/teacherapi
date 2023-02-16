package request

// Structure for "/api/retrievefornotifications" endpoint request body.
type ReceieveForNotificationsRequest struct {
	Teacher      string `json:"teacher" binding:"required"`
	Notification string `json:"notification" binding:"required"`
}
