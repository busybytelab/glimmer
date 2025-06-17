package app

import (
	"os"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/rs/zerolog/log"
)

// createUserAccount creates a new account for a verified user and returns the account ID
func createUserAccount(e *core.RecordConfirmVerificationRequestEvent, username string) (string, error) {
	// Check if account already exists
	var accounts []*core.Record
	err := e.App.RecordQuery("accounts").
		AndWhere(dbx.NewExp("owner = {:owner}", dbx.Params{"owner": e.Record.Id})).
		Limit(1).
		All(&accounts)
	if err != nil {
		return "", err
	}
	if len(accounts) > 0 {
		return accounts[0].Id, nil
	}

	// Account doesn't exist, create it
	collection, err := e.App.FindCollectionByNameOrId("accounts")
	if err != nil {
		log.Error().Err(err).Msg("Failed to find accounts collection")
		return "", err
	}

	account := core.NewRecord(collection)
	account.Set("name", username+"'s account")
	account.Set("owner", e.Record.Id)

	if err := e.App.Save(account); err != nil {
		log.Error().Err(err).Msg("Failed to create account for verified user")
		return "", err
	}
	log.Info().Str("email", e.Record.Email()).Str("account", account.Id).Msg("Created account for verified user")
	return account.Id, nil
}

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

	// Register hook for account creation after email verification
	app.pb.OnRecordConfirmVerificationRequest("users").BindFunc(func(e *core.RecordConfirmVerificationRequestEvent) error {
		// Log the verification status before processing
		log.Info().
			Str("email", e.Record.Email()).
			Bool("verified_before", e.Record.GetBool("verified")).
			Msg("Processing email verification")

		// Get the user's email
		email := e.Record.Email()

		// Extract username from email (everything before @)
		username := strings.Split(email, "@")[0]

		// Set the user's name to the username
		e.Record.Set("name", username)

		// Create an account for the user
		_, err := createUserAccount(e, username)
		if err != nil {
			return err
		}

		// Try to fetch fresh record to ensure we have latest state
		collection, err := e.App.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		record, err := e.App.FindRecordById(collection.Id, e.Record.Id)
		if err != nil {
			return err
		}

		// Set verified flag
		record.Set("verified", true)

		// Save the updated user record with verified flag
		if err := e.App.Save(record); err != nil {
			log.Error().Err(err).Msg("Failed to update user record")
			return err
		}

		// Log the verification status after processing
		log.Info().
			Str("email", record.Email()).
			Bool("verified_after", record.GetBool("verified")).
			Msg("Completed email verification")

		return nil
	})
}
