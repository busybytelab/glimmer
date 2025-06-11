package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
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
					"id": "text3534813612",
					"name": "session_name",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "text2219694135",
					"name": "topic_name",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number4218096251",
					"name": "total_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number3430046044",
					"name": "answered_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number119232996",
					"name": "wrong_answers_count",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number2285871944",
					"name": "total_score",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "text2285871945",
					"name": "session_status",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "date2285871946",
					"name": "last_answer_time",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "date"
				},
				{
					"hidden": false,
					"id": "text2285871947",
					"name": "learner_id",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "number2285871948",
					"name": "approved_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number2285871949",
					"name": "edited_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"hidden": false,
					"id": "number2285871950",
					"name": "not_reviewed_items",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				}
			],
			"id": "pbc_practice_session_stats",
			"indexes": [],
			"name": "practice_session_stats",
			"system": false,
			"type": "view",
			"createRule": null,
			"deleteRule": null,
			"listRule": "@request.auth.id ?= @collection.accounts.owner && @collection.accounts.id ?= account",
			"updateRule": null,
			"viewRule": "@request.auth.id ?= @collection.accounts.owner && @collection.accounts.id ?= account",
			"viewQuery": "WITH session_stats AS (SELECT ps.id, ps.name as session_name, ps.status as session_status, ps.learner as learner_id, pt.name as topic_name, COUNT(pi.id) as total_items, COUNT(pr.id) as answered_items, SUM(CASE WHEN pr.is_correct = 0 THEN 1 ELSE 0 END) as wrong_answers_count, SUM(COALESCE(pr.score, 0)) as total_score, MAX(pr.submitted_at) as last_answer_time, ps.account as account, SUM(CASE WHEN pi.status = 'Approved' THEN 1 ELSE 0 END) as approved_items, SUM(CASE WHEN pi.status = 'Edited' THEN 1 ELSE 0 END) as edited_items, SUM(CASE WHEN pi.status = 'Generated' OR pi.status IS NULL THEN 1 ELSE 0 END) as not_reviewed_items FROM practice_sessions ps LEFT JOIN practice_topics pt ON ps.practice_topic = pt.id LEFT JOIN practice_items pi ON pi.id IN (SELECT value FROM json_each(ps.practice_items)) LEFT JOIN practice_results pr ON pr.practice_session = ps.id AND pr.practice_item = pi.id GROUP BY ps.id, ps.name, ps.status, ps.learner, pt.name, ps.account) SELECT id, session_name, topic_name, total_items, answered_items, wrong_answers_count, total_score, session_status, last_answer_time, learner_id, account, approved_items, edited_items, not_reviewed_items FROM session_stats ORDER BY CASE WHEN session_status != 'completed' THEN 1 WHEN wrong_answers_count > 0 THEN 2 ELSE 3 END, last_answer_time DESC"
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_practice_session_stats")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
