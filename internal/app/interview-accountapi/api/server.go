package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/commandhandlers"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/eventhandlers"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/executors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/processors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/queries"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/giantswarm/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
)

func Configure() {
	viper.AutomaticEnv()
	viper.SetDefault("MessageVisibilityTimeout", 60)

	db := connectToDatabase()

	migrateDatabase(db.DB)

	settings.ApplicationClientId, settings.ApplicationClientSecret = getApplicationCredentials()

	executors.Configure(db)

	processors.Configure()
	eventhandlers.Configure()
	commandhandlers.Configure()
	queries.Configure()

	setupRoutes()
}

func connectToDatabase() *sqlx.DB {
	var db *sqlx.DB
	var err error

	_ = retry.Do(func() error {
		db, err = sqlx.Connect("sqlite3", ":memory:")
		if err != nil {
			return err
		}

		return nil

	}, retry.MaxTries(10), retry.Sleep(time.Duration(200*time.Millisecond)))

	if err != nil {
		panic(err)
	}

	return db
}

func migrateDatabase(db *sql.DB) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := "internal/app/interview-accountapi/api/migrations/"
	if strings.HasSuffix(dir,"internal/app/interview-accountapi") {
		path = "api/migrations/"
	}

	n, err := migrate.Exec(db, "postgres", &migrate.FileMigrationSource{Dir: path}, migrate.Up)
	if err != nil {
		panic(fmt.Sprintf("could not migrate database, error: %v", err))
	}

	log.Infof("applied %d database migrations!\n", n)
}

func getApplicationCredentials() (string, string) {
	return viper.GetString(settings.ServiceName + "-credentials-client-id"), viper.GetString(settings.ServiceName + "-credentials-client-secret")
}

func setupRoutes() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(setupGinLogger())

	http.HandleFunc("/", router.ServeHTTP)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	v1 := router.Group("/v1")
	v1.GET("/health", HandleGetHealth)

	accounts := v1.Group("/organisation/accounts").Use(gin.Logger())
	{
		accounts.GET("/:id", WithUserContext(HandleGetAccountById))
		accounts.DELETE("/:id", WithUserContext(HandleDeleteAccount))
		accounts.GET("", WithUserContext(HandleListAccounts))
		accounts.POST("", WithUserContext(HandleCreateAccount))
	}

}

// setupGinLogger adds a Gin middleware that logs on endpoints other than the health check
func setupGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/v1/health" {
			gin.Logger()(c)
		}
	}
}

func StartServer(ch <-chan bool, startedSignal chan bool) {
	port := settings.ServerPort
	address := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: address, Handler: nil}
	go func() {
		log.Infof("Server started on %s", address)
		startedSignal <- true
		<-ch
		log.Info("Shutting down")
		_ = server.Shutdown(context.Background())
	}()
	log.Info(fmt.Sprintf("listening on localhost:%d", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

type connectionString struct {
	host     string
	port     int
	user     string
	password string
	database string
	sslMode  string
}

func newConnectionString(host string, user string, password string, database string, port int, sslMode string) connectionString {
	return connectionString{
		host:     host,
		user:     user,
		password: password,
		database: database,
		port:     port,
		sslMode:  sslMode,
	}
}

func (c connectionString) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		c.host, c.port, c.user, c.password, c.database, c.sslMode)

}

func getOrDefaultString(property string, defaultValue string) string {
	if viper.IsSet(property) {
		return viper.GetString(property)
	} else {
		return defaultValue
	}
}

func getOrDefaultInt(property string, defaultValue int) int {
	if viper.IsSet(property) {
		return viper.GetInt(property)
	} else {
		return defaultValue
	}
}
