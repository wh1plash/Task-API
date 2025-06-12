# CRUD operations with Task

## Test assignment
Implement an HTTP server in Go that provides CRUD operations on the Task entity. The goal is to test knowledge of Go (modules, packages, mutex, interfaces), ability to work with HTTP, JSON. Implementation of HTTP requests should be implemented through the Gin package.
The full description of the task is contained in "tz.pdf"

Before run application you must create .env file in root directory and set parameters

### Example of .env file
```
PG_HOST="db"
PG_PORT=5432
PG_USER="postgres"
PG_PASS="postgres"
PG_DB_NAME="Task_CRUD"
LISTENADDR="0.0.0.0:3000"
```
Set __PG_HOST="db"__ if you vant to run application in Docker with docker-compose or __PG_HOST="localhost"__ if you vant to run locally.
Make sure you Postgres instance up and running before start application

### Use commands from Makefile to run application.
run with docker-compose
Linux
```
make docker
```
Windows
```
docker compose up --build -d
```
run app locally
Linux
```
make run
```
Windows
```
go run main.go
```
run seed for add 20 Tasks with random status
Linux
```
make seed
```
Windows
```
go run scripts/seed.go
```
## Add Task example JSON body
POST http://localhost:3000/tasks
```
{
    "title": "Task4",
    "description": "4",
    "status": "open"
}
```
## Get Tasks with pagination and query-param
GET http://localhost:3000/tasks?page=1&page_size=150&status=open
