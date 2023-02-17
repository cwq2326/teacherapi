package middlewares

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Registers middleware to router.
func RegisterDatabaseMiddleware(router *gin.Engine, db *sql.DB) {
	router.Use(func(c *gin.Context) {
		DatabaseMiddleware(c, db)
	})
}

func DatabaseMiddleware(c *gin.Context, db *sql.DB) {
	c.Set("db", db)
	c.Next()
}
