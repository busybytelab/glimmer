package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_instructors")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id = account.owner.id || @request.auth.id = user.id || @collection.learners.account.id = account.id",
			"viewRule": "@request.auth.id = user.id || @request.auth.id = account.owner.id || @collection.learners.account.id = account.id"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_instructors")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"listRule": "@request.auth.id = account.owner.id || @request.auth.id = user.id",
			"viewRule": "@request.auth.id = user.id || @request.auth.id = account.owner.id"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
