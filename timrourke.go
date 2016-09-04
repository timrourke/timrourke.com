package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timrourke/timrourke.com/db"
	"log"
	"os"
	"time"
)

func main() {
	// Load dot env files
	err := godotenv.Load()
	if err != nil {
		logError(err)
		panic(err)
	}

	// Initialize DB connection
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	mysqlDBName := os.Getenv("MYSQL_DBNAME")
	DB, err := db.ConnectToDB(mysqlUser, mysqlPass, mysqlDBName)
	if err != nil {
		logError(err)
		panic(err)
	}

	fmt.Println(DB)

	r := initRouter()

	r.Run(":8000")
}

// Logs an error to stderr, prefixed with a timestamp
func logError(err error) {
	now := time.Now().Format(time.RFC3339Nano)
	log.Println(fmt.Sprintf("%v:", now), err)
}

// Initialize gin-gonic routes
func initRouter() *gin.Engine {
	fmt.Println("environment", os.Getenv("ENV"))
	r := gin.Default()

	r.GET("/ping", getPing)

	return r
}

// Ping route handler for testing uptime
func getPing(c *gin.Context) {
	c.String(200, "pong")
}
