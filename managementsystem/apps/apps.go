package apps

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/phannita016/management/apis/collect/authorize"
	"github.com/phannita016/management/apis/collect/management"
	"github.com/phannita016/management/apis/middleware"
	"github.com/phannita016/management/apis/serve"
	"github.com/phannita016/management/driver"
	"github.com/phannita016/management/store"
	"github.com/phannita016/management/store/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

type Apps struct {
	// addr listen http server.
	Addr string

	// context integrate config collector.
	Ctx context.Context

	// http-server
	server *Server

	// driver mongo client
	mongoDB *mongo.Client

	// internal log handle.
	Logger *slog.Logger

	// secret authorize
	// sign the token with a secret key.
	Secret []byte

	// skipper defines a function to skip middleware.
	Skippers []string

	// cache store
	Cache store.Cache[string]

	// middleware control handler
	Middleware *middleware.Middleware
}

// AppsServer running application
func AppsServer(app *Apps, opts ...Option) func(c context.Context) error {
	// apply apps default config.
	app.ApplyApps()

	// apply config
	for _, opt := range opts {
		opt(app)
	}

	// srv instant serve echo server
	// serve input middleware argument security.
	var e = serve.New(middleware.Logger(app.Logger), app.Middleware.Authorization(), app.Middleware.AuthorizeWithToken)

	// assign apps-server
	if app.server == nil {
		app.server = NewServer(app.Addr, e, app.Logger)
	}

	// connect driver mongoDB
	db, err := app.ConnectDriver()
	if err != nil {
		return func(c context.Context) error {
			return err
		}
	}
	app.mongoDB = db

	// handler route apis
	app.HandleFunc(e)

	// listen appx collector server.
	// binding port addr
	if err := app.server.Run(app.Ctx); err != nil {
		app.Logger.ErrorContext(app.Ctx, "Server", slog.Any("err", err))
	}

	// return function shutdown server
	return app.server.Stop
}

// HandleFunc function route handler
func (apps *Apps) HandleFunc(e *echo.Echo) {
	// route authorized handler
	authorize.New(e, apps.Secret, apps.Cache)

	// route management handler
	g := e.Group("/management")
	management.New(g, apps.mongoDB)
}

// ApplyApps set default value of struct.
func (apps *Apps) ApplyApps() {
	if apps.Logger == nil {
		apps.Logger = slog.Default()
	}
	if apps.Ctx == nil {
		apps.Ctx = context.Background()
	}
	if apps.Skippers == nil {
		apps.Skippers = []string{
			"/authorize/login",
			"/authorize/refresh",
		}
	}
	if apps.Cache == nil {
		apps.Cache = cache.NewMemCache[string](cache.MemTTL(time.Hour * 24))
	}
	if apps.Middleware == nil {
		apps.Middleware = middleware.NewMiddleware(apps.Secret, apps.Cache, apps.Skippers)
	}
}

// new mongo client
func (apps *Apps) ConnectDriver() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. This is okay in production if environment variables are set.")
		return nil, err
	}

	uri := os.Getenv("MONGODB_URI")
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	client, err := driver.NewMongoClient(driver.MongoDriver{
		Hostname: uri,
		Username: username,
		Password: password,
		PoolSize: 20,
	})
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("err", err))
		return nil, err
	}

	slog.Info("Connection Database-Server...")
	return client, nil
}
