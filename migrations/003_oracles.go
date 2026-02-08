package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		oracles := core.NewBaseCollection("oracles")

		oracles.Fields.Add(&core.TextField{Name: "name", Required: true})
		oracles.Fields.Add(&core.TextField{Name: "description"})
		oracles.Fields.Add(&core.TextField{Name: "birth_issue"})
		oracles.Fields.Add(&core.TextField{Name: "verification_issue"})
		oracles.Fields.Add(&core.TextField{Name: "github_repo"})
		oracles.Fields.Add(&core.TextField{Name: "owner_wallet"})
		oracles.Fields.Add(&core.TextField{Name: "bot_wallet"})
		oracles.Fields.Add(&core.BoolField{Name: "wallet_verified"})
		oracles.Fields.Add(&core.BoolField{Name: "approved"})
		oracles.Fields.Add(&core.BoolField{Name: "claimed"})
		oracles.Fields.Add(&core.NumberField{Name: "karma"})
		oracles.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		oracles.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		oracles.Indexes = append(oracles.Indexes,
			"CREATE UNIQUE INDEX idx_oracles_birth ON oracles(birth_issue)",
		)

		listRule := ""
		viewRule := ""
		oracles.ListRule = &listRule
		oracles.ViewRule = &viewRule

		return app.Save(oracles)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("oracles")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
