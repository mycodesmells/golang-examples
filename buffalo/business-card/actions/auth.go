package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	c.Session().Set("token", user.AccessToken)

	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}

func IsAuth() buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			t := c.Session().Get("token")
			if t == nil {
				c.Session().Set("login_redirect_to", c.Request().URL.String())
				return c.Redirect(http.StatusFound, "/auth/github")
			}
			return h(c)
		}
	}
}

func AuthLogout(c buffalo.Context) error {
	c.Session().Delete("token")
	return c.Redirect(http.StatusFound, "/")
}
