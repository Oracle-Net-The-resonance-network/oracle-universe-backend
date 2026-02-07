package hooks

import (
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
