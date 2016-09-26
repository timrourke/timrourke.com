package main

import (
	"errors"
	"flag"
	"fmt"
	//	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
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
	"net/http"
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

	// Expose DB to models for relationship resolutions
	model.DB = DB

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

// Add CORS headers to fin requwsts
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}

// Add CORS headers to api2go requests
func Api2goCorsMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH, OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// Initialize gin-gonic routes
func initRouter(DB *sqlx.DB) *gin.Engine {
	r := gin.Default()

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("http://localhost:8000"),
		gingonic.New(r),
	)

	api.UseMiddleware(Api2goCorsMiddleware)

	userStorage := storage.NewUserStorage(DB)
	api.AddResource(model.User{}, resource.UserResource{
		UserStorage: userStorage,
	})

	postStorage := storage.NewPostStorage(DB)
	api.AddResource(model.Post{}, resource.PostResource{
		PostStorage: postStorage,
	})

	r.GET("/ping", getPing)

	r.Use(static.Serve("/", static.LocalFile("./frontend/html", true)))
	// r.Use(static.Serve("/", static.LocalFile("./hugo/public", true)))

	r.Use(static.Serve("/admin", static.LocalFile("./admin/dist", true)))

	r.GET("/admin/*wildcard", func(c *gin.Context) {
		path := fmt.Sprintf("%v", c.Request.URL)
		path = "./admin/dist/index.html"
		fmt.Println("trying to load path: ", path)
		http.ServeFile(c.Writer, c.Request, path)
	})

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
