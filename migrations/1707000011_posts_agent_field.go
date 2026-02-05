package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		agents, err := app.FindCollectionByNameOrId("agents")
		if err != nil {
			return err
		}

		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}

		// Add agent relation field to posts
		// Posts can be authored by either a human (author field) or an agent (agent field)
		posts.Fields.Add(&core.RelationField{
			Name:         "agent",
			CollectionId: agents.Id,
			MaxSelect:    1,
		})

		return app.Save(posts)
	}, func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return nil
		}

		// Remove agent field on rollback
		posts.Fields.RemoveByName("agent")
		return app.Save(posts)
	})
}
