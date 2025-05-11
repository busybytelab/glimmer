package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_accounts")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id != \"\" && (@collection.learners.user ?= @request.auth.id && @collection.learners.account ?= id || @collection.instructors.user ?= @request.auth.id && @collection.instructors.account ?= id)",
			"viewRule": "@request.auth.id != \"\" && (@collection.learners.user ?= @request.auth.id && @collection.learners.account ?= id || @collection.instructors.user ?= @request.auth.id && @collection.instructors.account ?= id)"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_accounts")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id = owner.id",
			"viewRule": "@request.auth.id = owner.id"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
