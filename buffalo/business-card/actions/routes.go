package actions

import "github.com/gobuffalo/buffalo"

// RoutesHandler is a default handler to serve up
// a routes page.
func RoutesHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("routes.html", "old_application.html"))
}
