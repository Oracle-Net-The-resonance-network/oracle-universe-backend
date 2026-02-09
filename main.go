package main

import (
	"log"
	"net/http"
	"os"

	"oracle-universe/hooks"
	_ "oracle-universe/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

// Injected at build time via -ldflags
var (
	Version   = "dev"
	BuildDate = "unknown"
	CommitSHA = "unknown"
)

func main() {
	app := pocketbase.New()

	// Enable auto migrations
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// Register record lifecycle hooks (API routes handled by CF Workers)
	hooks.RegisterHooks(app)

	// GET /api/version — build info
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/version", func(e *core.RequestEvent) error {
			return e.JSON(http.StatusOK, map[string]string{
				"version":   Version,
				"buildDate": BuildDate,
				"commit":    CommitSHA,
			})
		})
		return se.Next()
	})

	// Start the server
	if err := app.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
