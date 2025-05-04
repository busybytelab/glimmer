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
			"createRule": "@request.auth.id != null",
			"deleteRule": "@request.auth.id = instructor.user.id",
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
					"id": "name_column",
					"maxSize": 2000,
					"name": "name",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "status_column",
					"maxSize": 2000,
					"name": "status",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "assigned_at_column",
					"name": "assigned_at",
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
					"id": "completed_at_column",
					"name": "completed_at",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "date"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "generation_prompt_column",
					"maxSize": 2000,
					"name": "generation_prompt",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "learner_column",
					"maxSize": 2000,
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
					"id": "practice_topic_column",
					"maxSize": 2000,
					"name": "practice_topic",
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
					"maxSize": 2000,
					"name": "account",
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
					"id": "practice_items_column",
					"maxSize": 2000,
					"name": "practice_items",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "relation",
					"collectionId": "pbc_%s",
					"cascadeDelete": false,
					"maxSelect": null,
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
			"listRule": "@request.auth.id != null",
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id != null",
			"viewRule": "@request.auth.id != null"
		}`, domain.CollectionLearners, domain.CollectionPracticeTopics, domain.CollectionAccounts, domain.CollectionPracticeItems, domain.CollectionPracticeSessions, domain.CollectionPracticeSessions)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
