package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"task/api"
	"task/store"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	mustLoadEnvVariables()
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), port, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	db, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatal("error to connect to Posgres database", "error", err.Error())
		return
	}

	if err := db.Init(); err != nil {
		log.Fatal("error to create tables", err)
		return
	}

	gin.SetMode(gin.DebugMode)
	app := gin.Default()
	tastHandler := api.NewTaskHandler(db)
	app.POST("/tasks", tastHandler.HandlePostTask)
	app.GET("/tasks", tastHandler.HandleGetTasks)
	app.GET("/tasks/:id", tastHandler.HandleGetTask)
	app.PUT("/tasks/:id", tastHandler.HandlePutTask)
	app.DELETE("/tasks/:id", tastHandler.HandleDeleteTask)
	// app.GET("/validate", middleware.RequireAuth, controllers.Validate)
	// app.POST("/upload", middleware.RequireAuth, controllers.UploadHandler)
	// app.POST("/upload/csv", middleware.RequireAuth, controllers.UploadAndConcert)

	log.Fatal("Error to star HTTP server", app.Run(os.Getenv("LISTENADDR")))

}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
