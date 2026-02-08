package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		posts := core.NewBaseCollection("posts")

		posts.Fields.Add(&core.TextField{Name: "title", Required: true})
		posts.Fields.Add(&core.TextField{Name: "content"})
		posts.Fields.Add(&core.TextField{Name: "author_wallet", Required: true})
		posts.Fields.Add(&core.TextField{Name: "oracle_birth_issue"})
		posts.Fields.Add(&core.NumberField{Name: "upvotes"})
		posts.Fields.Add(&core.NumberField{Name: "downvotes"})
		posts.Fields.Add(&core.NumberField{Name: "score"})
		posts.Fields.Add(&core.TextField{Name: "tags"})
		posts.Fields.Add(&core.TextField{Name: "siwe_message"})
		posts.Fields.Add(&core.TextField{Name: "siwe_signature"})
		posts.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		posts.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		listRule := ""
		viewRule := ""
		posts.ListRule = &listRule
		posts.ViewRule = &viewRule

		return app.Save(posts)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("posts")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
