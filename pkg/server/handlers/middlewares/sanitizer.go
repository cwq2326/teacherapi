package middlewares

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"

	"govtech/pkg/utilities/messages"
)

// Registers middleware to router.
func RegisterSantizerMiddleware(router *gin.Engine) {
	router.Use(sanitizeRequestBody)
	router.Use(sanitizeQueryParams)
}

// This function sanitize the request body if content type is JSON.
func sanitizeRequestBody(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")

	// Check request content type is json.
	if contentType == "application/json" {
		// Sanitize the request body fields.
		var fields map[string]interface{}
		if err := c.ShouldBindJSON(&fields); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": messages.MESSAGE_BAD_REQUEST})
			return
		}
		for k, v := range fields {
			fields[k] = html.EscapeString(fmt.Sprintf("%v", v))
		}

		// Replace the request body with the sanitized fields.
		c.Set("json", fields)
	}
}

// This function sanitize the query parameters.
func sanitizeQueryParams(c *gin.Context) {
	// Sanitize the query parameters.
	query := c.Request.URL.Query()
	for k, v := range query {
		for i := range v {
			query.Set(k, html.EscapeString(v[i]))
		}
	}
	c.Request.URL.RawQuery = query.Encode()
}
