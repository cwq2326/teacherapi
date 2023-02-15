package controllers

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"govtech/pkg/utilities/messages"
)

func RegisterCommonStudentsEndpoint(r *gin.Engine, db *sql.DB) {
	r.GET("/api/commonstudents", func(c *gin.Context) {
		CommonStudents(c, db)
	})
}

/*
This function handles a GET request to the "/api/commonstudents" endpoint.
It returns all students common to a given list of teachers.
*/
func CommonStudents(c *gin.Context, db *sql.DB) {
	teachers := c.QueryArray("teacher")

	// Return error reponse if no "teacher" query parameter is given.
	if len(teachers) == 0 {
		var parameters = []string{"teacher"}

		c.JSON(http.StatusBadRequest, gin.H{"message": messages.MissingQueryParamsMessage(parameters)})
		return
	}

	var query string
	var student string
	var students []string

	// Build query string to get students registered to all teachers in the list.
	for i, v := range teachers {
		if i > 0 {
			query += " INTERSECT "
		}
		query += `SELECT student
				  FROM teaches
				  WHERE teacher = "` + v + `"`
	}

	// Query DB to get all students.
	result, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
		return
	}

	for result.Next() {
		err := result.Scan(&student)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": messages.MESSAGE_DATABASE_ERROR})
			return
		}
		students = append(students, student)
	}
	sort.Strings(students)
	
	c.JSON(http.StatusOK, gin.H{"students": students})
}
