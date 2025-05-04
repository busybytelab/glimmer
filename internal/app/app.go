package app

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/rs/zerolog/log"
)

// Application is the core glimmer service, with all required dependencies.
type Application struct {
	pb           *pocketbase.PocketBase
	quit         chan os.Signal
	shutdownOnce sync.Once
	publicFS     fs.FS
	config       *Config
}

// create a new application instance with the provided filesystem for static files.
// if publicFS is nil, it will fall back to using the local pb_public directory.
func New(publicFS fs.FS) *Application {
	log.Info().Msg("Creating new application instance...")
	return &Application{
		pb:       pocketbase.New(),
		quit:     make(chan os.Signal, 1),
		publicFS: publicFS,
		config:   EnvConfig(),
	}
}

func (app *Application) Initialize() error {
	log.Info().Msg("Initializing application...")

	app.setupMigrations()

	app.setupRoutes()

	app.setupCollectionsAndHooks()

	app.setupGracefulShutdown()

	return nil
}

// configures database migrations
func (a *Application) setupMigrations() {
	log.Info().Bool("autoMigrate", a.config.DB.AutoMigrate).Msg("Setting up migrations...")
	migratecmd.MustRegister(a.pb, a.pb.RootCmd, migratecmd.Config{
		Automigrate: a.config.DB.AutoMigrate,
	})
}

// configures the HTTP routes for the application
func (app *Application) setupRoutes() {
	app.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// serve static files - use the embedded filesystem if provided, otherwise use local directory
		staticFS := app.publicFS
		if staticFS == nil {
			log.Info().Msg("Using physical pb_public directory for static files")
			staticFS = os.DirFS("./pb_public")
		} else {
			log.Info().Msg("Using embedded static files")
		}

		e.Router.GET("/{path...}", apis.Static(staticFS, false))

		// add custom API endpoints here
		// e.Router.GET("/api/custom", func(c *core.RequestEvent) error { return nil })

		// must call e.Next() to continue the serve chain
		return e.Next()
	})
}

// register PocketBase collections and hooks
func (app *Application) setupCollectionsAndHooks() {
	// TODO: register collections and hooks
	log.Warn().Msg("Collections and hooks pending implementation")
}

// configure signal handling for graceful shutdown
func (app *Application) setupGracefulShutdown() {
	// register for SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(app.quit, os.Interrupt, syscall.SIGTERM)

	// hook into PocketBase's termination event
	app.pb.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		log.Info().Msg("PocketBase termination event triggered")
		app.Shutdown()
		return nil
	})
}

// starts the PocketBase server
func (app *Application) Start() error {
	log.Info().Msg("Starting application...")

	if app.pb == nil {
		log.Error().Msg("Cannot start: PocketBase instance not initialized")
		return errors.New("pocketbase instance not initialized")
	}

	// listen for termination signals
	go func() {
		<-app.quit
		log.Info().Msg("Shutdown signal received")
		app.Shutdown()
	}()

	log.Info().Msg("Starting PocketBase server...")
	return app.pb.Start()
}

// gracefully shuts down the application
func (app *Application) Shutdown() {
	app.shutdownOnce.Do(func() {
		log.Info().Msg("Shutting down gracefully...")

		// create a context with a timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// any cleanup operations can go here
		_ = ctx // Use ctx in actual shutdown operations when needed

		// close the quit channel to indicate shutdown is complete
		select {
		case <-app.quit: // channel already received signal
			// do nothing
		default:
			close(app.quit)
		}

		log.Info().Msg("Shutdown complete")
	})
}
