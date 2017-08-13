package actions

import "github.com/gobuffalo/buffalo"

// ContactHandler is a default handler to serve up
// a contact page.
func ContactHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("contact.html"))
}
