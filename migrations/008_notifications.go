package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		notifications := core.NewBaseCollection("notifications")

		notifications.Fields.Add(&core.TextField{Name: "recipient_wallet", Required: true})
		notifications.Fields.Add(&core.TextField{Name: "actor_wallet", Required: true})
		notifications.Fields.Add(&core.TextField{Name: "type", Required: true})
		notifications.Fields.Add(&core.TextField{Name: "message", Required: true})
		notifications.Fields.Add(&core.TextField{Name: "post_id"})
		notifications.Fields.Add(&core.TextField{Name: "comment_id"})
		notifications.Fields.Add(&core.NumberField{Name: "count"})
		notifications.Fields.Add(&core.BoolField{Name: "read"})
		notifications.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		notifications.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		notifications.Indexes = append(notifications.Indexes,
			"CREATE INDEX idx_notifications_recipient ON notifications(recipient_wallet, read, created)",
			"CREATE INDEX idx_notifications_vote_batch ON notifications(recipient_wallet, type, post_id, read)",
		)

		listRule := ""
		viewRule := ""
		notifications.ListRule = &listRule
		notifications.ViewRule = &viewRule

		return app.Save(notifications)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("notifications")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
