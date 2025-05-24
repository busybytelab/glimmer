package app

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/busybytelab.com/glimmer/internal/llm"
	chatRoutePkg "github.com/busybytelab.com/glimmer/internal/route/chat"
	llmRoutePkg "github.com/busybytelab.com/glimmer/internal/route/llm"
	practiceRoutePkg "github.com/busybytelab.com/glimmer/internal/route/practice"
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
	llmService   llm.Service
	chatService  llm.ChatService
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

func (app *Application) Initialize() error {
	log.Debug().Msg("Initializing application...")
	app.setupMigrations()
	app.setupLLMService()
	app.setupRoutes()
	app.setupCollectionsAndHooks()
	app.setupCommands()
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

// configures custom commands
func (app *Application) setupCommands() {
	log.Trace().Msg("Setting up custom commands...")
	setupSeedCommand(app.pb)
	log.Trace().Msg("Custom commands setup completed")
}

// configures the HTTP routes for the application
func (app *Application) setupRoutes() {
	llmRoutes := llmRoutePkg.New(app.llmService)
	practiceRoute := practiceRoutePkg.NewPracticeSessionRoute(app.llmService)
	answerRoute := practiceRoutePkg.NewAnswerRoute()
	chatRoutes := chatRoutePkg.New(app.chatService)

	app.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// API routes - register these first for priority

		// LLM chat API endpoint for common chat requests
		e.Router.POST("/api/glimmer/v1/llm/chat", llmRoutes.HandleChatRequest).Bind(apis.RequireAuth())
		e.Router.GET("/api/glimmer/v1/llm/info", llmRoutes.HandleInfoRequest).Bind(apis.RequireAuth())
		e.Router.POST("/api/glimmer/v1/practice/session", practiceRoute.HandleCreatePracticeSession).Bind(apis.RequireAuth())
		e.Router.POST("/api/glimmer/v1/practice/evaluate-answer", answerRoute.HandleEvaluateAnswer).Bind(apis.RequireAuth())

		// Chat API endpoints
		e.Router.POST("/api/glimmer/v1/chat", chatRoutes.HandleChatRequest).Bind(apis.RequireAuth())

		// Create a custom handler for static files that ensures correct MIME types
		staticHandler := func(c *core.RequestEvent) error {
			// Get requested path
			reqPath := c.Request.URL.Path

			// Skip API routes explicitly
			if strings.HasPrefix(reqPath, "/api/") {
				return c.NotFoundError("API route not found", nil)
			}

			// JavaScript modules need specific MIME types
			if strings.HasSuffix(reqPath, ".js") {
				// Set the handler to set the Content-Type header
				fs := http.FileServer(http.FS(app.embedFs))
				w := c.Response

				// Set the correct MIME type for JavaScript modules
				w.Header().Set("Content-Type", "application/javascript")

				// Serve the file
				fs.ServeHTTP(w, c.Request)
				return nil
			}

			// For other static files, use the default handler
			return apis.Static(app.embedFs, true)(c)
		}

		// Handle all non-API static requests with the custom handler
		e.Router.GET("/{path...}", staticHandler)

		// must call e.Next() to continue the serve chain
		return e.Next()
	})
}

// register PocketBase collections and hooks
func (app *Application) setupCollectionsAndHooks() {
}

// initialize the LLM service
func (app *Application) setupLLMService() {
	llmConfig := llm.LoadConfig()

	// Setup with PocketBase app for cache storage if needed
	app.llmService = llm.AppService(llmConfig, app.pb)

	log.Info().Msg("LLM service initialized")
	// Initialize the chat service with PocketBase app and LLM service
	app.chatService = llm.NewChatService(app.pb, app.llmService)

	log.Info().Msg("Chat service initialized")
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
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Note: PocketBase handles database connections internally
		// No need to manually close them as it's done by PocketBase's shutdown process

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
