package app

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cast"
)

// configureAppSettings sets up application settings from environment variables
func (app *Application) configureAppSettings() {
	app.pb.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		// Configure app metadata
		e.App.Settings().Meta.AppName = os.Getenv("APP_NAME")
		e.App.Settings().Meta.AppURL = os.Getenv("APP_URL")
		e.App.Settings().Meta.SenderName = os.Getenv("SENDER_NAME")
		e.App.Settings().Meta.SenderAddress = os.Getenv("SENDER_ADDRESS")

		// Configure SMTP settings from environment variables
		e.App.Settings().SMTP.Enabled = cast.ToBool(os.Getenv("SMTP_ENABLED"))
		e.App.Settings().SMTP.Host = os.Getenv("SMTP_HOST")
		e.App.Settings().SMTP.Port = cast.ToInt(os.Getenv("SMTP_PORT"))
		e.App.Settings().SMTP.Username = os.Getenv("SMTP_USERNAME")
		e.App.Settings().SMTP.Password = os.Getenv("SMTP_PASSWORD")
		e.App.Settings().SMTP.AuthMethod = os.Getenv("SMTP_AUTH_METHOD")
		e.App.Settings().SMTP.TLS = cast.ToBool(os.Getenv("SMTP_TLS"))
		e.App.Settings().SMTP.LocalName = os.Getenv("SMTP_LOCAL_NAME")

		// Validate and persist the changes
		return e.App.Save(e.App.Settings())
	})

	// Prevent SMTP settings from being changed through the admin UI
	app.pb.OnSettingsUpdateRequest().BindFunc(func(e *core.SettingsUpdateRequestEvent) error {
		if e.OldSettings.SMTP.Enabled != e.NewSettings.SMTP.Enabled ||
			e.OldSettings.SMTP.Host != e.NewSettings.SMTP.Host ||
			e.OldSettings.SMTP.Port != e.NewSettings.SMTP.Port ||
			e.OldSettings.SMTP.Username != e.NewSettings.SMTP.Username ||
			e.OldSettings.SMTP.Password != e.NewSettings.SMTP.Password ||
			e.OldSettings.SMTP.AuthMethod != e.NewSettings.SMTP.AuthMethod ||
			e.OldSettings.SMTP.TLS != e.NewSettings.SMTP.TLS ||
			e.OldSettings.SMTP.LocalName != e.NewSettings.SMTP.LocalName {
			return e.ForbiddenError("Cannot change the SMTP settings", nil)
		}
		return e.Next()
	})
}
