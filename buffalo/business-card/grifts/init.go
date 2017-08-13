package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/mycodesmells/golang-examples/buffalo/business-card/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
