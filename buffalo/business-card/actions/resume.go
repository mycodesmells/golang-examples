package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mycodesmells/golang-examples/buffalo/business-card/models"
	"github.com/pkg/errors"
)

// ResumeHandler is a default handler to serve up
// a resume page.
func ResumeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	experience, err := buildExperience(tx)
	if err != nil {
		return errors.Wrap(err, "failed to load experience data")
	}
	c.Set("experience", experience)

	skillset, err := buildSkillset(tx)
	if err != nil {
		return errors.Wrap(err, "failed to load skills date")
	}
	c.Set("skillset", skillset)

	return c.Render(200, r.HTML("resume.html"))
}

func buildExperience(tx *pop.Connection) ([]models.Experience, error) {
	var experience []models.Experience
	if err := tx.All(&experience); err != nil {
		return nil, err
	}
	return experience, nil
}

func buildSkillset(tx *pop.Connection) (map[string][]models.Skill, error) {
	var skills []models.Skill
	skillset := make(map[string][]models.Skill)
	if err := tx.All(&skills); err != nil {
		return nil, err
	}
	for _, s := range skills {
		category := s.Category.String
		if category == "" {
			category = "Other"
		}
		skillset[category] = append(skillset[category], s)
	}
	return skillset, nil
}
