package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"
	"github.com/timrourke/timrourke.com/db"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/resource"
	"github.com/timrourke/timrourke.com/storage"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	var serve, test bool

	set := flag.NewFlagSet("timrourke", flag.ExitOnError)
	set.BoolVar(&serve, "serve", false, "Serve timrourke.com.")
	set.BoolVar(&test, "test", false, "Run tests for timrourke.com.")
	set.Parse(os.Args[1:])

	if test && serve {
		err := errors.New("You cannot both test and run the app simultaneously.")
		logError(err)
		panic(err)
	}

	if test {
		runTests()
		os.Exit(0)
	} else if serve {
		run()
	} else {
		err := errors.New("You must supply a command line argument to use this binary.")
		logError(err)
		panic(err)
	}
}

// Logs an error to stderr, prefixed with a timestamp
func logError(err error) {
	now := time.Now().Format(time.RFC3339Nano)
	log.Println(fmt.Sprintf("%v:", now), err)
}

// Load dot env files
func loadDotEnv(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		logError(err)
		panic(err)
	}
}

// Run the application
func run() {
	// Load development or production dotenv
	loadDotEnv(".env")

	// Initialize DB connection
	DB := initDB()

	// Initialize routes
	r := initRouter(DB)

	r.Run(":8000")
}

// Connect to database
func initDB() *sqlx.DB {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	mysqlDBName := os.Getenv("MYSQL_DBNAME")

	DB, err := db.ConnectToDB(mysqlUser, mysqlPass, mysqlDBName)
	if err != nil {
		logError(err)
		panic(err)
	}

	return DB
}

// Initialize gin-gonic routes
func initRouter(DB *sqlx.DB) *gin.Engine {
	r := gin.Default()

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("http://localhost:8000"),
		gingonic.New(r),
	)

	userStorage := storage.NewUserStorage(DB)
	api.AddResource(model.User{}, resource.UserResource{
		UserStorage: userStorage,
	})

	r.GET("/ping", getPing)

	return r
}

// Ping route handler for testing uptime
func getPing(c *gin.Context) {
	c.String(200, "pong")
}

// Run database migration for testing environment
func migrateDB(migrationDirection string) ([]byte, error) {
	cmd := exec.Command("migrate",
		"-url",
		fmt.Sprintf("mysql://%s:%s@/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_DBNAME")),
		"-path",
		"./migrations",
		migrationDirection)

	return cmd.CombinedOutput()
}

func runTests() {
	// Load test-specific dotenv
	loadDotEnv(".env.test")

	// Apply migrations for test db
	out, err := migrateDB("up")
	if err != nil {
		fmt.Print(string(out))
		logError(err)
		panic(err)
	}
	fmt.Print(string(out))

	// Run godog test suites
	cmd := exec.Command("godog")
	cmd.Args[0] = "./features"
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		logError(err)
		panic(err)
	}
	fmt.Print(string(out))

	// Undo migrations for test db
	out, err = migrateDB("down")
	if err != nil {
		fmt.Print(string(out))
		logError(err)
		panic(err)
	}
	fmt.Print(string(out))
}
