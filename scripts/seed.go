package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"task/store"
	"task/store/fixtures"

	"github.com/joho/godotenv"
)

func init() {
	mustLoadEnvVariables()
}

type status int

const (
	Open status = iota
	InProgress
	Done
)

var statusStrings = [...]string{
	"open",
	"in_progress",
	"done",
}

func (s status) String() string {
	return statusStrings[s]
}

func main() {
	var (
		port, _ = strconv.Atoi(os.Getenv("PG_PORT"))
		connStr = fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", port, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	)
	db, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatal("error to connect to Posgres database: ", err.Error())
		return
	}

	if err := db.Init(); err != nil {
		log.Fatal("error to create tables", err)
		return
	}

	for i := range 20 {
		title := fmt.Sprintf("random task title %d", i)
		description := fmt.Sprintf("description_%d", i)
		statusIndex := rand.Intn(len(statusStrings))
		status := status(statusIndex).String()
		fixtures.AddTasks(db, title, description, status)
	}
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
