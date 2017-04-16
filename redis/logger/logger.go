package logger

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func FromContext(ectx echo.Context) *logrus.Entry {
	sessionID := ectx.Get("sessionID").(string)
	return logrus.WithField("sessionID", sessionID)
}
