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
	"github.com/rs/zerolog/log"
)

// Application represents the core glimmer service, with all required dependencies.
type Application struct {
	pb           *pocketbase.PocketBase
	quit         chan os.Signal
	shutdownOnce sync.Once
	publicFS     fs.FS
}

// New creates a new application instance with the provided filesystem for static files.
// If publicFS is nil, it will fall back to using the local pb_public directory.
func New(publicFS fs.FS) *Application {
	log.Info().Msg("Creating new application instance...")
	return &Application{
		pb:       pocketbase.New(),
		quit:     make(chan os.Signal, 1),
		publicFS: publicFS,
	}
}

func (app *Application) Initialize() error {
	log.Info().Msg("Initializing application...")

	// TODO(v0.1): Implement configuration loading from config.yaml - see issue #XX

	// Setup HTTP routes
	app.setupRoutes()

	// Register collections and hooks
	app.setupCollectionsAndHooks()

	// Setup graceful shutdown
	app.setupGracefulShutdown()

	return nil
}

// setupRoutes configures the HTTP routes for the application
func (app *Application) setupRoutes() {
	app.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Serve static files - use the embedded filesystem if provided, otherwise use local directory
		staticFS := app.publicFS
		if staticFS == nil {
			log.Info().Msg("Using physical pb_public directory for static files")
			staticFS = os.DirFS("./pb_public")
		} else {
			log.Info().Msg("Using embedded static files")
		}

		e.Router.GET("/{path...}", apis.Static(staticFS, false))

		// Add custom API endpoints here
		// e.Router.GET("/api/custom", func(c *core.RequestEvent) error { return nil })

		// Must call e.Next() to continue the serve chain
		return e.Next()
	})
}

// setupCollectionsAndHooks registers PocketBase collections and hooks
func (app *Application) setupCollectionsAndHooks() {
	// TODO: Register collections and hooks
	log.Warn().Msg("Collections and hooks pending implementation")
}

// setupGracefulShutdown configures signal handling for graceful shutdown
func (app *Application) setupGracefulShutdown() {
	// Register for SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(app.quit, os.Interrupt, syscall.SIGTERM)

	// Hook into PocketBase's termination event
	app.pb.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		log.Info().Msg("PocketBase termination event triggered")
		app.Shutdown()
		return nil
	})
}

// Start runs the application (starts the PocketBase server)
func (app *Application) Start() error {
	log.Info().Msg("Starting application...")

	if app.pb == nil {
		log.Error().Msg("Cannot start: PocketBase instance not initialized")
		return errors.New("pocketbase instance not initialized")
	}

	// Start a goroutine to listen for termination signals
	go func() {
		<-app.quit
		log.Info().Msg("Shutdown signal received")
		app.Shutdown()
	}()

	log.Info().Msg("Starting PocketBase server...")
	return app.pb.Start()
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() {
	app.shutdownOnce.Do(func() {
		log.Info().Msg("Shutting down gracefully...")

		// Create a context with a timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Any cleanup operations can go here
		_ = ctx // Use ctx in actual shutdown operations when needed

		// Close the quit channel to indicate shutdown is complete
		select {
		case <-app.quit: // Channel already received signal
			// Do nothing
		default:
			close(app.quit)
		}

		log.Info().Msg("Shutdown complete")
	})
}
