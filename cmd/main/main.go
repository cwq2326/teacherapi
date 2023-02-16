package main

import (
	"fmt"
	"path/filepath"
	"os"

	"github.com/joho/godotenv"

	"govtech/pkg/server/databases"
	"govtech/pkg/server/handlers"
)

var dbConfig database.MySqlConfig
var routerConfig handlers.RouterConfig

func init() {
	// Init env
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		fmt.Println("Failed to load .env file")
	}

	dbConfig = database.MySqlConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
	}

	routerConfig = handlers.RouterConfig{
		Port: os.Getenv("ROUTER_PORT"),
		Host: os.Getenv("ROUTER_HOST"),
	}

	deployPort := os.Getenv("HTTP_PLATFORM_PORT")
	if (deployPort != "") {
		routerConfig.Port = deployPort
	}
}

func main() {
	// Init database.
	db := database.ConnectDB(&dbConfig)
	database.InitDB(db)
	defer database.DisconnectDB(db)

	// Init router.
	r := handlers.InitRouter()
	handlers.RegisterEndpoints(r, db)
	handlers.RegisterMiddlewares(r)
	handlers.RunRouter(r, &routerConfig)
}
