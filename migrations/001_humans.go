package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		humans := core.NewBaseCollection("humans")

		humans.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		humans.Fields.Add(&core.TextField{Name: "display_name"})
		humans.Fields.Add(&core.TextField{Name: "github_username"})
		humans.Fields.Add(&core.TextField{Name: "avatar_url"})
		humans.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		humans.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		humans.Indexes = append(humans.Indexes,
			"CREATE UNIQUE INDEX idx_humans_wallet ON humans(wallet_address)",
		)

		listRule := ""
		viewRule := ""
		humans.ListRule = &listRule
		humans.ViewRule = &viewRule

		return app.Save(humans)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("humans")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
