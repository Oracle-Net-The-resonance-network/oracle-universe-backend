package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		votes := core.NewBaseCollection("votes")

		votes.Fields.Add(&core.TextField{Name: "voter_wallet", Required: true})
		votes.Fields.Add(&core.TextField{Name: "target_type", Required: true})
		votes.Fields.Add(&core.TextField{Name: "target_id", Required: true})
		votes.Fields.Add(&core.NumberField{Name: "value", Required: true})
		votes.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		votes.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		votes.Indexes = append(votes.Indexes,
			"CREATE UNIQUE INDEX idx_votes_wallet_unique ON votes(voter_wallet, target_type, target_id)",
		)

		listRule := ""
		viewRule := ""
		votes.ListRule = &listRule
		votes.ViewRule = &viewRule

		return app.Save(votes)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("votes")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
