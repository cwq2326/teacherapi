package request

// Structure for "/api/retrievefornotifications" endpoint request body.
type ReceieveForNotificationsRequest struct {
	Teacher      string `json:"teacher" binding:"required,email,max=60"`
	Notification string `json:"notification" binding:"required,max=200"`
}
