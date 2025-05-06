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
			"deleteRule": null,
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
					"id": "key_column",
					"max": 0,
					"min": 0,
					"name": "key",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"convertURLs": false,
					"hidden": false,
					"id": "prompt_column",
					"max": 200000,
					"name": "prompt",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "editor"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "system_prompt_column",
					"max": 200000,
					"name": "system_prompt",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "editor"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "response_column",
					"max": 200000,
					"name": "response",
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
					"id": "model_name_column",
					"max": 20000,
					"name": "model_name",
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
					"id": "backend_column",
					"max": 20000,
					"name": "backend",
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
					"id": "total_tokens_column",
					"max": null,
					"min": 0,
					"name": "total_tokens",
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
					"id": "prompt_tokens_column",
					"max": null,
					"min": 0,
					"name": "prompt_tokens",
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
					"id": "completion_tokens_column",
					"max": null,
					"min": 0,
					"name": "completion_tokens",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "cost_column",
					"max": null,
					"min": 0,
					"name": "cost",
					"onlyInt": false,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "ttl_column",
					"max": 0,
					"min": 0,
					"name": "ttl",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "number"
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
			"indexes": [
				"CREATE UNIQUE INDEX `+"`"+`idx_key_unique`+"`"+` ON `+"`"+`%s`+"`"+` (`+"`"+`key`+"`"+`)"
			],
			"listRule": null,
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": null,
			"viewRule": null
		}`, domain.CollectionLLMResponses, domain.CollectionLLMResponses, domain.CollectionLLMResponses)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionLLMResponses)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
