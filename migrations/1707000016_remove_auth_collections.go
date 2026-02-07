package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// ============================================================
		// Convert humans from AuthCollection to BaseCollection
		// All auth is SIWE + custom JWT + admin token — PB auth fields unused
		// ============================================================

		// Drop old auth collection
		humansOld, err := app.FindCollectionByNameOrId("humans")
		if err == nil {
			if err := app.Delete(humansOld); err != nil {
				return err
			}
		}

		// Recreate as BaseCollection (no email/password/tokenKey)
		humans := core.NewBaseCollection("humans")
		humans.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		humans.Fields.Add(&core.TextField{Name: "display_name"})
		humans.Fields.Add(&core.TextField{Name: "github_username"})
		humans.Fields.Add(&core.TextField{Name: "avatar_url"})
		humans.Indexes = append(humans.Indexes, "CREATE UNIQUE INDEX idx_humans_wallet ON humans(wallet_address)")
		if err := app.Save(humans); err != nil {
			return err
		}

		// ============================================================
		// Convert agents from AuthCollection to BaseCollection
		// ============================================================

		agentsOld, err := app.FindCollectionByNameOrId("agents")
		if err == nil {
			if err := app.Delete(agentsOld); err != nil {
				return err
			}
		}

		agents := core.NewBaseCollection("agents")
		agents.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		agents.Fields.Add(&core.TextField{Name: "display_name"})
		agents.Fields.Add(&core.NumberField{Name: "reputation"})
		agents.Fields.Add(&core.BoolField{Name: "verified"})
		agents.Indexes = append(agents.Indexes, "CREATE UNIQUE INDEX idx_agents_wallet ON agents(wallet_address)")
		return app.Save(agents)

	}, func(app core.App) error {
		// Rollback: convert back to AuthCollections
		// (DB wipes on deploy anyway, but good practice)

		// Drop base agents
		agentsBase, err := app.FindCollectionByNameOrId("agents")
		if err == nil {
			app.Delete(agentsBase)
		}
		agents := core.NewAuthCollection("agents")
		agents.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		agents.Fields.Add(&core.TextField{Name: "display_name"})
		agents.Fields.Add(&core.NumberField{Name: "reputation"})
		agents.Fields.Add(&core.BoolField{Name: "verified"})
		agents.Indexes = append(agents.Indexes, "CREATE UNIQUE INDEX idx_agents_wallet ON agents(wallet_address)")
		if err := app.Save(agents); err != nil {
			return err
		}

		// Drop base humans
		humansBase, err := app.FindCollectionByNameOrId("humans")
		if err == nil {
			app.Delete(humansBase)
		}
		humans := core.NewAuthCollection("humans")
		humans.Fields.Add(&core.TextField{Name: "wallet_address", Required: true})
		humans.Fields.Add(&core.TextField{Name: "display_name"})
		humans.Fields.Add(&core.TextField{Name: "github_username"})
		humans.Fields.Add(&core.TextField{Name: "avatar_url"})
		humans.Indexes = append(humans.Indexes, "CREATE UNIQUE INDEX idx_humans_wallet ON humans(wallet_address)")
		return app.Save(humans)
	})
}
