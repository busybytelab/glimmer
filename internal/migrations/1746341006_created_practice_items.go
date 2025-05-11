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
					"id": "question_text_column",
					"max": 2000,
					"name": "question_text",
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
					"id": "question_type_column",
					"max": 2000,
					"name": "question_type",
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
					"id": "options_column",
					"max": 2000,
					"name": "options",
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
					"id": "correct_answer_column",
					"max": 2000,
					"name": "correct_answer",
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
					"id": "explanation_column",
					"max": 2000,
					"name": "explanation",
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
					"id": "explanation_for_incorrect_column",
					"max": 2000,
					"name": "explanation_for_incorrect",
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
					"id": "hints_column",
					"max": 2000,
					"name": "hints",
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
					"id": "difficulty_level_column",
					"max": 2000,
					"name": "difficulty_level",
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
					"max": 2000,
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
					"id": "tags_column",
					"max": 2000,
					"name": "tags",
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
					"id": "practice_topic_column",
					"max": 2000,
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
			"name": "%s",
			"system": false,
			"type": "base",
			"createRule": "@request.auth.id = account.owner",
			"deleteRule": "@request.auth.id = account.owner",
			"listRule": "@request.auth.id = account.owner || (@collection.learners.account ?= account && @collection.learners.user = @request.auth.id)",
			"updateRule": "@request.auth.id = account.owner",
			"viewRule": "@request.auth.id = account.owner || (@collection.learners.account ?= account && @collection.learners.user = @request.auth.id)"
		}`, domain.CollectionPracticeTopics, domain.CollectionAccounts, domain.CollectionPracticeItems, domain.CollectionPracticeItems)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
