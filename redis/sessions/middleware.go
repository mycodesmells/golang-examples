package sessions

import (
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

func Middleware(store Store) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			cookie, err := ectx.Cookie("sessionID")
			if err != nil {
				sessionID := uuid.NewV4().String()
				ectx.SetCookie(&http.Cookie{
					Name:  "sessionID",
					Value: sessionID,
				})
				ectx.Set("sessionID", sessionID)
				store.Set(sessionID, Session{})
				return hf(ectx)
			}

			sessionID := cookie.Value
			ectx.Set("sessionID", sessionID)
			return hf(ectx)
		}
	}
}
