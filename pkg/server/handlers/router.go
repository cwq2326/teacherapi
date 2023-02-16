package handlers

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	"govtech/pkg/controllers"
	"govtech/pkg/server/handlers/middlewares"
)

// Structure for configuration for the router.
type RouterConfig struct {
	Port string
	Host string
}

// Returns an instance of the router.
func InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}

// Runs the router as the assigned port and host.
func RunRouter(router *gin.Engine, config *RouterConfig) {
	connectionString := fmt.Sprintf("%s:%s", config.Host, config.Port)
	router.Run(connectionString)
}

// Register endpoints to the router.
func RegisterEndpoints(router *gin.Engine, db *sql.DB) {
	endpointRegistrations := []func(*gin.Engine, *sql.DB) {
		controllers.RegisterCommonStudentsEndpoint,
		controllers.RegisterRegisterEndpoint,
		controllers.RegisterRetrieveForNotificationEndpoint,
		controllers.RegisterSuspendEndpoint,
	}

	for _,v := range endpointRegistrations {
		v(router, db)
	}
}


// Register middlewares to the router.
func RegisterMiddlewares(router *gin.Engine) {
	middlewares := []func(*gin.Engine) {
		middlewares.RegisterSantizerMiddleware,
	}

	for _,v := range middlewares {
		v(router)
	}
}
