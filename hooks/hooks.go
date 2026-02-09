package hooks

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterHooks sets up record lifecycle hooks for Oracle Universe
// All auth is handled by the Elysia API (SIWE + custom JWT + admin token).
// Hooks only initialize default field values.
func RegisterHooks(app *pocketbase.PocketBase) {
	// ============================================================
	// AGENT WORLD - Sandbox for AI entities
	// ============================================================

	// Agents: Set defaults
	app.OnRecordCreateRequest("agents").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Record.Set("reputation", 0)
		e.Record.Set("verified", false)
		return e.Next()
	})

	// ============================================================
	// HUMAN WORLD - Verified wallet holders
	// ============================================================

	// Posts: Initialize vote counters
	// The API layer sets author/agent/oracle fields via admin token
	app.OnRecordCreateRequest("posts").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Record.Set("upvotes", 0)
		e.Record.Set("downvotes", 0)
		e.Record.Set("score", 0)
		return e.Next()
	})

	// Comments: Initialize vote counters
	// The API layer sets author field via admin token
	app.OnRecordCreateRequest("comments").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Record.Set("upvotes", 0)
		e.Record.Set("downvotes", 0)
		return e.Next()
	})

	// ============================================================
	// NOTIFICATIONS — auto-created when a comment is saved
	// ============================================================

	// After a comment is persisted, notify the post owner (human) and oracle (bot)
	app.OnRecordCreate("comments").BindFunc(func(e *core.RecordEvent) error {
		// Let the create finish first
		if err := e.Next(); err != nil {
			return err
		}

		// Fire-and-forget: notification failure must never block comment creation
		go func() {
			defer func() { recover() }() // swallow panics

			postID := e.Record.GetString("post")
			actorWallet := strings.ToLower(e.Record.GetString("author_wallet"))
			commentID := e.Record.Id
			if postID == "" || actorWallet == "" {
				return
			}

			post, err := app.FindRecordById("posts", postID)
			if err != nil {
				return
			}

			title := post.GetString("title")
			msg := fmt.Sprintf(`commented on "%s"`, title)
			birthIssue := post.GetString("oracle_birth_issue")

			// Resolve recipients
			var ownerWallet, botWallet string

			if birthIssue != "" {
				// Oracle post — look up oracle
				oracles, err := app.FindRecordsByFilter("oracles",
					fmt.Sprintf(`birth_issue="%s"`, birthIssue), "", 1, 0)
				if err == nil && len(oracles) > 0 {
					ownerWallet = strings.ToLower(oracles[0].GetString("owner_wallet"))
					botWallet = strings.ToLower(oracles[0].GetString("bot_wallet"))
				}
			}

			// Human post fallback
			if ownerWallet == "" {
				ownerWallet = strings.ToLower(post.GetString("author_wallet"))
			}

			notifCol, err := app.FindCollectionByNameOrId("notifications")
			if err != nil {
				return
			}

			// Helper: create notification with self-suppression (same wallet only)
			notify := func(recipient string) {
				if recipient == "" || recipient == actorWallet {
					return
				}

				n := core.NewRecord(notifCol)
				n.Set("recipient_wallet", recipient)
				n.Set("actor_wallet", actorWallet)
				n.Set("type", "comment")
				n.Set("message", msg)
				n.Set("post_id", postID)
				n.Set("comment_id", commentID)
				n.Set("read", false)
				app.Save(n)
			}

			// 1. Human notification
			notify(ownerWallet)
			// 2. Oracle bot notification (separate inbox)
			if botWallet != ownerWallet {
				notify(botWallet)
			}
		}()

		return nil
	})

	// Oracles: Set defaults (only if not already set by API)
	app.OnRecordCreateRequest("oracles").BindFunc(func(e *core.RecordRequestEvent) error {
		if !e.Record.GetBool("approved") {
			e.Record.Set("approved", false)
		}
		if e.Record.GetFloat("karma") == 0 {
			e.Record.Set("karma", 0)
		}
		return e.Next()
	})

	// ============================================================
	// NOTE: All API routes moved to Elysia (apps/api)
	// - /api/oracles, /api/posts, /api/feed → apps/api/routes/
	// - /api/humans/me, /api/agents/me → apps/api/routes/
	// - /skill.md, /docs, /openapi.json → apps/api/
	// ============================================================
}
