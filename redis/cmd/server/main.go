package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/mycodesmells/golang-examples/redis/logger"
	pagehit "github.com/mycodesmells/golang-examples/redis/pagehit-2"
	"github.com/mycodesmells/golang-examples/redis/sessions"
)

func main() {
	sessionsStore := sessions.NewRedisStore()
	hitStore := pagehit.NewRedisStore()

	e := echo.New()
	e.Use(sessions.Middleware(sessionsStore))
	e.Use(pagehit.Middleware(hitStore))

	e.GET("/", func(ectx echo.Context) error {
		log := logger.FromContext(ectx)

		sessionID := ectx.Get("sessionID").(string)
		s, err := sessionsStore.Get(sessionID)
		if err != nil {
			log.Errorf("err: %v", err)
		}

		log.Infof("Visits: %d", s.VisitCount)
		response := fmt.Sprintf("Hello World #%d\n", s.VisitCount)

		s.VisitCount = s.VisitCount + 1
		err = sessionsStore.Set(sessionID, s)
		if err != nil {
			log.Errorf("err: %v", err)
		}

		return ectx.String(http.StatusOK, response)
	})
	e.GET("/stats", pagehit.Handler(hitStore))

	e.Start(":5000")
}
