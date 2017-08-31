package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/mycodesmells/golang-examples/buffalo/business-card/models"
	"github.com/pkg/errors"
)

// ResumeHandler is a default handler to serve up
// a resume page.
func ResumeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	var experience []models.Experience
	if err := tx.All(&experience); err != nil {
		return errors.Wrap(err, "failed to load experience data")
	}
	c.Set("experience", experience)
	return c.Render(200, r.HTML("resume.html"))
}
