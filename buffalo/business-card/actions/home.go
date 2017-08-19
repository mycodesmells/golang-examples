package actions

import (
	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	tfn := c.Value("t").(func(string) (string, error))
	msg, _ := tfn("welcome_greeting")

	//c.Logger().Infof("data: %v", c.Data())
	c.Logger().Infof("data 1: %v", msg)
	c.Logger().Infof("data 2: %v", translate(c, "welcome_greeting"))
	return c.Render(200, r.HTML("home.html"))
}

func translate(ctx buffalo.Context, key string) string {
	tfn, ok := ctx.Value("t").(func(string) (string, error))
	if !ok {
		return key
	}
	msg, err := tfn(key)
	if err != nil {
		return key
	}

	return msg
}
