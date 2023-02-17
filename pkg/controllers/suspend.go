package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"govtech/pkg/models/request"
	"govtech/pkg/utilities/messages"
)

func RegisterSuspendEndpoint(r *gin.Engine) {
	r.POST("/api/suspend", Suspend)
}

/*
This function handles a POST request to the "/api/suspend" endpoint.
It suspends the specified student.
*/
func Suspend(c *gin.Context) {
	var request request.SuspendRequest
	db := c.MustGet("db").(*sql.DB)

	// Return error response if missing or invalid request body fields.
	if err := c.ShouldBindJSON(&request); err != nil {
		bindErr, paramErr := messages.GetErrorMessage(err)

		if bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": bindErr})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": paramErr})
			return
		}
	}

	// Update `suspended` field of specified student to 1 to indicate suspension.
	_, err := db.Query(`UPDATE students
						SET suspended = 1
						WHERE email = (?)`, request.Student)

	// Return error response if there is an error while querying the DB.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
		return
	}

	c.Status(http.StatusNoContent)
}
