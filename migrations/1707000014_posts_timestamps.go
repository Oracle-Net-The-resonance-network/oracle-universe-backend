package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}

		// Add created and updated autodate fields if they don't exist
		// These are used for sorting posts by date
		if posts.Fields.GetByName("created") == nil {
			posts.Fields.Add(&core.AutodateField{
				Name:     "created",
				OnCreate: true,
				OnUpdate: false,
			})
		}

		if posts.Fields.GetByName("updated") == nil {
			posts.Fields.Add(&core.AutodateField{
				Name:     "updated",
				OnCreate: true,
				OnUpdate: true,
			})
		}

		return app.Save(posts)
	}, func(app core.App) error {
		// Rollback: remove timestamp fields
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return nil
		}

		posts.Fields.RemoveByName("created")
		posts.Fields.RemoveByName("updated")

		return app.Save(posts)
	})
}
