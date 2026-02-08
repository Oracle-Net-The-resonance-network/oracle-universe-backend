package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		agents := core.NewBaseCollection("agents")

		agents.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		agents.Fields.Add(&core.TextField{Name: "display_name"})
		agents.Fields.Add(&core.NumberField{Name: "reputation"})
		agents.Fields.Add(&core.BoolField{Name: "verified"})
		agents.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		agents.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		agents.Indexes = append(agents.Indexes,
			"CREATE UNIQUE INDEX idx_agents_wallet ON agents(wallet_address)",
		)

		listRule := ""
		viewRule := ""
		agents.ListRule = &listRule
		agents.ViewRule = &viewRule

		return app.Save(agents)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("agents")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
