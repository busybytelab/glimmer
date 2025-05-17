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
			"deleteRule": "@request.auth.id = chat.user",
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
					"hidden": false,
					"id": "chat_column",
					"name": "chat",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation",
					"collectionId": "pbc_%s",
					"cascadeDelete": true,
					"maxSelect": 1
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "role_column",
					"max": 20,
					"min": 0,
					"name": "role",
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
					"id": "content_column",
					"max": 200000,
					"name": "content",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "editor"
				},
				{
					"hidden": false,
					"id": "usage_column",
					"name": "usage",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "order_column",
					"max": null,
					"min": 0,
					"name": "order",
					"presentable": false,
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
				"CREATE INDEX `+"`"+`idx_chat_order`+"`"+` ON `+"`"+`%s`+"`"+` (`+"`"+`chat`+"`"+`, `+"`"+`order`+"`"+`)",
				"CREATE INDEX `+"`"+`idx_chat_items_content`+"`"+` ON `+"`"+`%s`+"`"+` (`+"`"+`content`+"`"+`)"
			],
			"listRule": "@request.auth.id = chat.user",
			"name": "%s",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = chat.user",
			"viewRule": "@request.auth.id = chat.user"
		}`, domain.CollectionChats, domain.CollectionChatItems, domain.CollectionChatItems, domain.CollectionChatItems, domain.CollectionChatItems)

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(domain.CollectionChatItems)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
