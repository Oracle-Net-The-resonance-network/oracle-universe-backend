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

		comments := core.NewBaseCollection("comments")

		comments.Fields.Add(&core.TextField{Name: "content", Required: true})
		comments.Fields.Add(&core.RelationField{
			Name:         "post",
			CollectionId: posts.Id,
			Required:     true,
			MaxSelect:    1,
		})
		comments.Fields.Add(&core.TextField{Name: "author_wallet", Required: true})
		comments.Fields.Add(&core.TextField{Name: "parent_id"})
		comments.Fields.Add(&core.NumberField{Name: "upvotes"})
		comments.Fields.Add(&core.NumberField{Name: "downvotes"})
		comments.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
		comments.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})

		listRule := ""
		viewRule := ""
		comments.ListRule = &listRule
		comments.ViewRule = &viewRule

		return app.Save(comments)
	}, func(app core.App) error {
		col, _ := app.FindCollectionByNameOrId("comments")
		if col != nil {
			return app.Delete(col)
		}
		return nil
	})
}
