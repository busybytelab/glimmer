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
			"createRule": "@request.auth.id = practice_item.account.owner.id || @collection.learners.account.id = practice_item.account.id",
			"deleteRule": "@request.auth.id = practice_item.account.owner.id || @collection.learners.account.id = practice_item.account.id",
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
					"id": "answer_column",
					"max": 2000,
					"name": "answer",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "json"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "is_correct_column",
					"name": "is_correct",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "bool"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "score_column",
					"max": 100,
					"min": 0,
					"name": "score",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "feedback_column",
					"max": 2000,
					"name": "feedback",
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
					"id": "hint_level_reached_column",
					"max": 10,
					"min": 0,
					"name": "hint_level_reached",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "started_at_column",
					"name": "started_at",
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
					"id": "submitted_at_column",
					"name": "submitted_at",
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
					"id": "attempt_number_column",
					"max": 10,
					"min": 1,
					"name": "attempt_number",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "evaluation_details_column",
					"max": 2000,
					"name": "evaluation_details",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "practice_item_column",
					"max": 2000,
					"name": "practice_item",
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
					"id": "practice_session_column",
					"max": 2000,
					"name": "practice_session",
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
			"listRule": "@request.auth.id = practice_item.account.owner.id || @collection.learners.account.id = practice_item.account.id",
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = practice_item.account.owner.id || @collection.learners.account.id = practice_item.account.id",
			"viewRule": "@request.auth.id = practice_item.account.owner.id || @collection.learners.account.id = practice_item.account.id"
		}`, domain.CollectionPracticeItems, domain.CollectionLearners, domain.CollectionPracticeSessions, domain.CollectionPracticeResults, domain.CollectionPracticeResults)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeResults)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
