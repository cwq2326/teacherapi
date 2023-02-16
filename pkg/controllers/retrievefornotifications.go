package controllers

import (
	"database/sql"
	"net/http"
	"regexp"
	"sort"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"govtech/pkg/models/request"
	"govtech/pkg/utilities/messages"
	"govtech/pkg/utilities/patterns"
	"govtech/pkg/utilities/set"
)

func RegisterRetrieveForNotificationEndpoint(r *gin.Engine, db *sql.DB) {
	r.POST("/api/retrievefornotifications", func (c *gin.Context) {
		ReceiveForNotifications(c, db)
	})
}

/*
This function handles a POST request to the "/api/retrievefornotifications" endpoint.
It returns all students who can receive a notification from a teacher.
A student can receive a notification if he is not suspended and is registered to the teacher
or is mentioned in the notification.
*/

func ReceiveForNotifications(c *gin.Context, db *sql.DB) {
	var request request.ReceieveForNotificationsRequest

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

	// Check if notifications follow the regexp pattern.
	match := patterns.ValidatePattern(
		patterns.REGEX_PATTERN_NOTIFICATION,
		request.Notification,
	)

	// Return error if notification string does not follow the regexp pattern.
	if !match {
		var fields = []string{"notification"}

		c.JSON(http.StatusBadRequest, gin.H{"message": messages.InvalidParamsMessage(fields)})
		return
	}

	set := set.New[string]()
	var student string

	// Get all students registered under the teacher who are not suspended.
	result, err := db.Query(`SELECT student
							 FROM students 
							 INNER JOIN teaches
							 ON students.email = teaches.student
							 WHERE students.suspended = 0
							 AND teaches.teacher = (?)`, request.Teacher)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
		return
	}

	for result.Next() {
		err := result.Scan(&student)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
			return
		}
		set.Add(student)
	}

	// Get all tagged students in the notification who are not
	// suspended and are in the database.
	var count int
	regex := regexp.MustCompile(patterns.REGEX_PATTERN_EMAIL)
	taggedStudents := regex.FindAllString(request.Notification, -1)
	for _, v := range taggedStudents {
		err := db.QueryRow(`SELECT COUNT(*)
							FROM students
							WHERE suspended = 0
							AND email = (?)`, v).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
			return
		}

		// Add student to set if it exists.
		if count == 1 {
			set.Add(v)
		}
	}
	array := set.ToArray()
	sort.Strings(array)

	c.JSON(http.StatusOK, gin.H{"recipient": array})
}
