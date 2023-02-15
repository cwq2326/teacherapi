package handlers

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	"govtech/pkg/controllers"
)

type RouterConfig struct {
	Port string
	Host string
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func RunRouter(router *gin.Engine, config *RouterConfig) {
	connectionString := fmt.Sprintf("%s:%s", config.Host, config.Port)
	router.Run(connectionString)
}

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
