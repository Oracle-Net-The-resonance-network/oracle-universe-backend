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

		// Add verification_issue text field — stores the oracle-identity issue URL
		if oracles.Fields.GetByName("verification_issue") == nil {
			oracles.Fields.Add(&core.TextField{
				Name: "verification_issue",
			})
		}

		return app.Save(oracles)
	}, func(app core.App) error {
		oracles, err := app.FindCollectionByNameOrId("oracles")
		if err != nil {
			return nil
		}

		oracles.Fields.RemoveByName("verification_issue")

		return app.Save(oracles)
	})
}
