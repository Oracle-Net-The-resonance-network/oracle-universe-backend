package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// ============================================================
		// Wallet = Identity migration
		//
		// Replace PocketBase relation-based authorship with wallet addresses.
		// The wallet that signed is the identity — survives DB wipes.
		// ============================================================

		// --- POSTS ---
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}

		// Clear rules that reference old fields before removing them
		posts.DeleteRule = nil
		posts.UpdateRule = nil
		posts.ListRule = nil
		posts.ViewRule = nil

		// Remove old relation fields
		posts.Fields.RemoveByName("author") // was: human relation
		posts.Fields.RemoveByName("agent")  // was: agent relation
		posts.Fields.RemoveByName("oracle") // was: oracle relation

		// Add wallet-based fields
		posts.Fields.Add(&core.TextField{Name: "author_wallet", Required: true})
		posts.Fields.Add(&core.TextField{Name: "oracle_birth_issue"}) // stable oracle identifier

		if err := app.Save(posts); err != nil {
			return err
		}

		// --- COMMENTS ---
		comments, err := app.FindCollectionByNameOrId("comments")
		if err != nil {
			return err
		}

		// Clear rules that may reference old fields
		comments.DeleteRule = nil
		comments.UpdateRule = nil
		comments.ListRule = nil
		comments.ViewRule = nil

		// Remove old relation field
		comments.Fields.RemoveByName("author") // was: human relation

		// Add wallet-based field
		comments.Fields.Add(&core.TextField{Name: "author_wallet", Required: true})

		if err := app.Save(comments); err != nil {
			return err
		}

		// --- VOTES ---
		votes, err := app.FindCollectionByNameOrId("votes")
		if err != nil {
			return err
		}

		// Clear rules that may reference old fields
		votes.DeleteRule = nil
		votes.UpdateRule = nil
		votes.ListRule = nil
		votes.ViewRule = nil

		// Remove old relation field
		votes.Fields.RemoveByName("human") // was: human relation

		// Add wallet-based field
		votes.Fields.Add(&core.TextField{Name: "voter_wallet", Required: true})

		// Replace unique index (old used "human", new uses "voter_wallet")
		votes.Indexes = nil
		votes.Indexes = append(votes.Indexes, "CREATE UNIQUE INDEX idx_votes_wallet_unique ON votes(voter_wallet, target_type, target_id)")

		if err := app.Save(votes); err != nil {
			return err
		}

		// --- CONNECTIONS (drop - replaced by owner_wallet on oracles) ---
		connections, err := app.FindCollectionByNameOrId("connections")
		if err == nil {
			if err := app.Delete(connections); err != nil {
				return err
			}
		}

		// --- ORACLES ---
		oracles, err := app.FindCollectionByNameOrId("oracles")
		if err != nil {
			return err
		}

		// Remove old human relation
		oracles.Fields.RemoveByName("human") // was: human relation

		// Add owner_wallet (human owner's wallet)
		oracles.Fields.Add(&core.TextField{Name: "owner_wallet"})

		// Rename wallet_address → bot_wallet for clarity
		// (PB doesn't have rename, so remove + add)
		oracles.Fields.RemoveByName("wallet_address")
		oracles.Fields.Add(&core.TextField{Name: "bot_wallet"})

		return app.Save(oracles)

	}, func(app core.App) error {
		// Rollback: DB wipes on deploy anyway, but good practice
		// Restore relation fields (would need collection IDs, skip for now)
		return nil
	})
}
