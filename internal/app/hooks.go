package app

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/rs/zerolog/log"
)

// setupUserHooks registers hooks for the users collection
func (app *Application) setupUserHooks() {
	// Register hook for automatic email verification after user creation
	app.pb.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		// Check if auto-verification is enabled via environment variable
		if os.Getenv("AUTO_SEND_VERIFICATION_EMAIL") == "true" {
			// Send verification email
			if err := mails.SendRecordVerification(app.pb, e.Record); err != nil {
				log.Error().Err(err).Msg("Failed to send verification email")
				return err
			}

			log.Info().Str("email", e.Record.Email()).Msg("Verification email sent automatically")
		}

		return nil
	})
}
