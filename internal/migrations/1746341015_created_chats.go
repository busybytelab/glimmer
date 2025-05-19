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
			"createRule": "@request.auth.id != \"\"",
			"deleteRule": "@request.auth.id = user",
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
					"id": "user_column",
					"max": 2000,
					"name": "user",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "relation",
					"collectionId": "_pb_users_auth_",
					"cascadeDelete": true,
					"maxSelect": 1
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "label_column",
					"max": 100,
					"min": 0,
					"name": "label",
					"pattern": "",
					"presentable": true,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "system_prompt_column",
					"max": 2000,
					"name": "system_prompt",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "editor"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "model_column",
					"max": 100,
					"name": "model",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "total_tokens_column",
					"max": null,
					"min": 0,
					"name": "total_tokens",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "total_cost_column",
					"max": null,
					"min": 0,
					"name": "total_cost",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "archived_column",
					"name": "archived",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "bool"
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
			"listRule": "@request.auth.id = user",
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = user",
			"viewRule": "@request.auth.id = user"
		}`, domain.CollectionChats, domain.CollectionChats)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionChats)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
