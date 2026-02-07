package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add created/updated autodate fields to all collections that need them.
		// Posts already have these (migration 14), but comments, votes, oracles,
		// humans, and agents do not.

		for _, colName := range []string{"comments", "votes", "oracles", "humans", "agents"} {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				continue
			}

			if col.Fields.GetByName("created") == nil {
				col.Fields.Add(&core.AutodateField{
					Name:     "created",
					OnCreate: true,
					OnUpdate: false,
				})
			}

			if col.Fields.GetByName("updated") == nil {
				col.Fields.Add(&core.AutodateField{
					Name:     "updated",
					OnCreate: true,
					OnUpdate: true,
				})
			}

			if err := app.Save(col); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		for _, colName := range []string{"comments", "votes", "oracles", "humans", "agents"} {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				continue
			}
			col.Fields.RemoveByName("created")
			col.Fields.RemoveByName("updated")
			app.Save(col)
		}
		return nil
	})
}
