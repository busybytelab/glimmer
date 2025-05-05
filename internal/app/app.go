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

	"github.com/busybytelab.com/glimmer/internal/handler"
	"github.com/busybytelab.com/glimmer/internal/llm"
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
	embedFs      fs.FS
	config       *Config
	llmService   *llm.Service
}

// create a new application instance with the provided filesystem for static files.
// if embedFs is nil, it will fall back to using the local ui/dist directory.
func New(embedFs fs.FS) *Application {
	log.Info().Msg("Creating new application instance...")
	pb := pocketbase.New()
	return &Application{
		pb:      pb,
		quit:    make(chan os.Signal, 1),
		embedFs: embedFs,
		config:  EnvConfig(),
	}
}

// PB returns the PocketBase instance
func (app *Application) PB() *pocketbase.PocketBase {
	return app.pb
}

// LLM returns the LLM service
func (app *Application) LLM() *llm.Service {
	return app.llmService
}

func (app *Application) Initialize() error {
	log.Info().Msg("Initializing application...")

	app.setupMigrations()

	app.setupRoutes()

	app.setupCollectionsAndHooks()

	app.setupGracefulShutdown()

	// Initialize LLM service
	app.setupLLMService()

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
		// Add LLM chat API endpoint
		e.Router.POST("/api/llm/chat", handler.HandleChat(app.llmService))

		// Serve UI static files
		e.Router.GET("/{path...}", apis.Static(app.embedFs, true))

		// must call e.Next() to continue the serve chain
		return e.Next()
	})
}

// register PocketBase collections and hooks
func (app *Application) setupCollectionsAndHooks() {
	// TODO: register collections and hooks
	log.Warn().Msg("Collections and hooks pending implementation")
}

// initialize the LLM service
func (app *Application) setupLLMService() {
	// Load LLM configuration
	llmConfig := llm.LoadConfig()

	// Setup with PocketBase app for cache storage if needed
	app.llmService = llm.AppService(llmConfig, app.pb)

	log.Info().Msg("LLM service initialized")
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
