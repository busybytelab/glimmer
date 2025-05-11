package migrations

import (
	"encoding/json"
	"fmt"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := fmt.Sprintf(`{
			"createRule": null,
			"deleteRule": "@request.auth.id = account.owner.id || @collection.instructors.user.id = @request.auth.id",
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "id_column",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "earned_at_column",
					"name": "earned_at",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "date"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "achievement_definition_column",
					"max": 2000,
					"name": "achievement_definition",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "relation",

						"collectionId": "pbc_%s",
						"cascadeDelete": false,
						"maxSelect": 1,
						"minSelect": 1
					
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "learner_column",
					"max": 2000,
					"name": "learner",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "relation",

						"collectionId": "pbc_%s",
						"cascadeDelete": true,
						"maxSelect": 1,
						"minSelect": 1
					
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "account_column",
					"max": 2000,
					"name": "account",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "relation",

						"collectionId": "pbc_%s",
						"cascadeDelete": true,
						"maxSelect": 1,
						"minSelect": 1
					
				},
				{
					"hidden": false,
					"id": "created_column",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "updated_column",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_%s",
			"indexes": [],
			"listRule": "@request.auth.id = account.owner.id || @collection.instructors.user.id = @request.auth.id || @collection.learners.account.id = account.id",
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = account.owner.id || @collection.instructors.user.id = @request.auth.id",
			"viewRule": "@request.auth.id = account.owner.id || @collection.instructors.user.id = @request.auth.id || @collection.learners.account.id = account.id"
		}`, domain.CollectionAchievementDefinitions, domain.CollectionLearners, domain.CollectionAccounts, domain.CollectionEarnedAchievements, domain.CollectionEarnedAchievements)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionEarnedAchievements)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
