package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// ============================================================
		// Set public read rules on collections that the API reads
		// Without these, all reads require superuser auth.
		// ============================================================

		// Posts: public read, anyone can list/view
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}
		publicRule := ""
		posts.ListRule = &publicRule
		posts.ViewRule = &publicRule
		if err := app.Save(posts); err != nil {
			return err
		}

		// Comments: public read
		comments, err := app.FindCollectionByNameOrId("comments")
		if err != nil {
			return err
		}
		comments.ListRule = &publicRule
		comments.ViewRule = &publicRule
		if err := app.Save(comments); err != nil {
			return err
		}

		// Oracles: public read
		oracles, err := app.FindCollectionByNameOrId("oracles")
		if err != nil {
			return err
		}
		oracles.ListRule = &publicRule
		oracles.ViewRule = &publicRule
		if err := app.Save(oracles); err != nil {
			return err
		}

		// Humans: public read
		humans, err := app.FindCollectionByNameOrId("humans")
		if err != nil {
			return err
		}
		humans.ListRule = &publicRule
		humans.ViewRule = &publicRule
		if err := app.Save(humans); err != nil {
			return err
		}

		// Agents: public read
		agents, err := app.FindCollectionByNameOrId("agents")
		if err != nil {
			return err
		}
		agents.ListRule = &publicRule
		agents.ViewRule = &publicRule
		if err := app.Save(agents); err != nil {
			return err
		}

		// Votes: public read
		votes, err := app.FindCollectionByNameOrId("votes")
		if err != nil {
			return err
		}
		votes.ListRule = &publicRule
		votes.ViewRule = &publicRule
		return app.Save(votes)

	}, func(app core.App) error {
		// Rollback: set all rules back to nil (superuser-only)
		for _, colName := range []string{"posts", "comments", "oracles", "humans", "agents", "votes"} {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				continue
			}
			col.ListRule = nil
			col.ViewRule = nil
			app.Save(col)
		}
		return nil
	})
}
