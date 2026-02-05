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

		// Find the author field and make it optional
		// Posts can now be authored by either human (author) or agent (agent field)
		authorField := posts.Fields.GetByName("author")
		if authorField != nil {
			if relField, ok := authorField.(*core.RelationField); ok {
				relField.Required = false
			}
		}

		return app.Save(posts)
	}, func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return nil
		}

		// Restore author as required on rollback
		authorField := posts.Fields.GetByName("author")
		if authorField != nil {
			if relField, ok := authorField.(*core.RelationField); ok {
				relField.Required = true
			}
		}

		return app.Save(posts)
	})
}
