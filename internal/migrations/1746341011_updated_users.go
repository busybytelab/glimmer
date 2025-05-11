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

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewRule": "id = @request.auth.id || (@collection.learners.user ?= id && @collection.learners.account ?= @collection.instructors.account && @collection.instructors.user ?= @request.auth.id)",
			"listRule": "id = @request.auth.id || (@collection.learners.user ?= id && @collection.learners.account ?= @collection.instructors.account && @collection.instructors.user ?= @request.auth.id)",
			"updateRule": "id = @request.auth.id || (@collection.learners.user ?= id && @collection.learners.account ?= @collection.instructors.account && @collection.instructors.user ?= @request.auth.id)",
			"createRule": "id = @request.auth.id || @collection.accounts.owner = @request.auth.id",
			"authAlert": {
				"enabled": false
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
