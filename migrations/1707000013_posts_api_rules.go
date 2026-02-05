package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}

		// Allow public read access to posts
		// Anyone can list and view posts (feed)
		posts.ListRule = types.Pointer("")
		posts.ViewRule = types.Pointer("")

		// Create requires auth (API handles validation)
		posts.CreateRule = types.Pointer("@request.auth.id != ''")

		// Update and delete only by the author (human) or agent
		// This allows post owners to edit/delete their own posts
		posts.UpdateRule = types.Pointer("author = @request.auth.id || agent = @request.auth.id")
		posts.DeleteRule = types.Pointer("author = @request.auth.id || agent = @request.auth.id")

		return app.Save(posts)
	}, func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return nil
		}

		// Restore to superuser only on rollback
		posts.ListRule = nil
		posts.ViewRule = nil
		posts.CreateRule = nil
		posts.UpdateRule = nil
		posts.DeleteRule = nil

		return app.Save(posts)
	})
}
