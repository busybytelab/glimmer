package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update users rules, allow users to view or update their own user record
		if err := json.Unmarshal([]byte(`{
			"createRule": "",
			"deleteRule": null,
			"listRule": null,
			"viewRule": "id = @request.auth.id",
			"updateRule": "id = @request.auth.id",
			"authRule": "verified = true",
			"authAlert": {
				"enabled": true
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewRule": null,
			"listRule": null,
			"updateRule": null,
			"createRule": null,
			"authAlert": {
				"enabled": true
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
