package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"govtech/pkg/utilities/messages"
	"govtech/pkg/models/request"
)

func RegisterRegisterEndpoint(r *gin.Engine, db *sql.DB) {
	r.POST("/api/register", func(c *gin.Context) { 
		Register(c, db)
	})
}

/*
This function handles a POST request to the "/api/register" endpoint.
It can either register a list of students to a teacher or
a list of teachers to a student.
*/
func Register(c *gin.Context, db *sql.DB) {
	var request request.RegisterRequest

	// Return error response if missing or invalid request body fields.
	if err := c.ShouldBindJSON(&request); err != nil {
		bindErr, paramErr := messages.GetErrorMessage(err)

		if bindErr != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": 123})
			return
		} else if paramErr != "" {
			c.JSON(http.StatusBadGateway, gin.H{"message": 123})
			return
		}
	}

	haveTeacher := request.Teacher != ""
	haveStudents := len(request.Students) > 0

	// Check for valid pair of teacher and students field.
	if (haveTeacher && !haveStudents) || (!haveTeacher && haveStudents) {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.MissingValidPairMessage("teacher", "students")})
		return
	}

	haveStudent := request.Student != ""
	haveTeachers := len(request.Teachers) > 0

	// Check for valid pair of student and teachers field.
	if (haveStudent && !haveTeachers) || (!haveStudent && haveTeachers) {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.MissingValidPairMessage("student", "teachers")})
		return
	}

	// Try adding list of students into teacher table first.
	canAddToTeacher := haveTeacher && haveStudents

	if canAddToTeacher {
		err := insertIntoDB(request.Teacher, request.Students, 1, db)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
		}
	}

	// Add list of teachers into student table if valid fields are provided.
	canAddToStudent := haveStudent && haveTeachers

	if canAddToStudent {
		err := insertIntoDB(request.Student, request.Teachers, 2, db)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
		}
	}

	c.Status(http.StatusNoContent)
}

/*
This function queries the DB based on the target, list and action argument provided.
*/
func insertIntoDB(target string, list []string, action int, db *sql.DB) error {
	// Action denotes which query to be perform to the DB.

	// action = 1
	// Insert list of students into teacher table.

	// action = 2
	// Insert list of teachers into student table.

	switch action {
	case 1:
		_, err := db.Query(`INSERT IGNORE INTO teachers
							VALUES (?)`, target)

		if err != nil {
			return err
		}

		for _, v := range list {
			_, err := db.Query(`INSERT IGNORE INTO students
							VALUES (?, 0)`, v)

			if err != nil {
				return err
			}

			_, err = db.Query(`INSERT IGNORE INTO teaches
								VALUES (?, ?)`, target, v)
			if err != nil {
				return err
			}
		}
		return nil
	case 2:
		_, err := db.Query(`INSERT IGNORE INTO students
							VALUES (?, 0)`, target)

		if err != nil {
			return err
		}

		for _, v := range list {
			_, err := db.Query(`INSERT IGNORE INTO teachers
							VALUES (?)`, v)

			if err != nil {
				return err
			}

			_, err = db.Query(`INSERT IGNORE INTO teaches
								VALUES (?, ?)`, v, target)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}
