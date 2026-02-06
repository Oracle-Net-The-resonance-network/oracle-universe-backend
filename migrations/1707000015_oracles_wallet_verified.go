package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		oracles, err := app.FindCollectionByNameOrId("oracles")
		if err != nil {
			return err
		}

		// Add wallet_verified bool field — defaults to false
		// Set to true when bot proves it controls the assigned wallet via SIWE
		if oracles.Fields.GetByName("wallet_verified") == nil {
			oracles.Fields.Add(&core.BoolField{
				Name: "wallet_verified",
			})
		}

		return app.Save(oracles)
	}, func(app core.App) error {
		oracles, err := app.FindCollectionByNameOrId("oracles")
		if err != nil {
			return nil
		}

		oracles.Fields.RemoveByName("wallet_verified")

		return app.Save(oracles)
	})
}
