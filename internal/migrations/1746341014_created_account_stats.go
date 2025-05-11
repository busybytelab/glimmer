package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text3208210256",
					"max": 0,
					"min": 0,
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
					"id": "json3534813612",
					"maxSize": 1,
					"name": "account_name",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json2219694135",
					"maxSize": 1,
					"name": "total_learners",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json2792794425",
					"maxSize": 1,
					"name": "total_instructors",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json4218096251",
					"maxSize": 1,
					"name": "total_practice_topics",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json3430046044",
					"maxSize": 1,
					"name": "total_practice_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json119232996",
					"maxSize": 1,
					"name": "total_practice_results",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json2285871944",
					"maxSize": 1,
					"name": "avg_practice_results_per_learner",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				}
			],
			"id": "pbc_account_stats",
			"indexes": [],
			"listRule": "@request.auth.id != \"\" && @collection.instructors.user ?= @request.auth.id && @collection.instructors.account ?= id",
			"name": "account_stats",
			"system": false,
			"type": "view",
			"updateRule": null,
			"viewQuery": "WITH account_stats AS (\n    SELECT \n        a.id as account_id,\n        a.name as account_name,\n        COUNT(DISTINCT l.id) as total_learners,\n        COUNT(DISTINCT i.id) as total_instructors,\n        COUNT(DISTINCT pt.id) as total_practice_topics,\n        COUNT(DISTINCT pi.id) as total_practice_items,\n        COUNT(DISTINCT pr.id) as total_practice_results\n    FROM accounts a\n    LEFT JOIN learners l ON l.account = a.id\n    LEFT JOIN instructors i ON i.account = a.id\n    LEFT JOIN practice_topics pt ON pt.account = a.id\n    LEFT JOIN practice_items pi ON pi.practice_topic = pt.id\n    LEFT JOIN practice_results pr ON pr.practice_item = pi.id\n    GROUP BY a.id, a.name\n)\nSELECT \n    account_id as id,\n    account_name,\n    total_learners,\n    total_instructors,\n    total_practice_topics,\n    total_practice_items,\n    total_practice_results,\n    (CASE \n        WHEN total_learners > 0 \n        THEN (ROUND((CAST(total_practice_results AS FLOAT) / total_learners), 2))\n        ELSE 0 \n    END) as avg_practice_results_per_learner\nFROM account_stats\nORDER BY account_name;",
			"viewRule": "@request.auth.id != \"\" && @collection.instructors.user ?= @request.auth.id && @collection.instructors.account ?= id"
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_account_stats")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
