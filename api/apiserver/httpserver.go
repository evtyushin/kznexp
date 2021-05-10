package apiserver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"kznexp/api/handlers"
	"kznexp/service"
)

type Api struct {
	http   *echo.Echo
	listen string
}

func New(s *service.Service, listen string, tasks chan *service.Item, logger *zap.Logger) *Api {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	setRouters(e, s, tasks, logger)

	return &Api{http: e, listen: listen}
}

func (a *Api) Start() error {
	return a.http.Start(a.listen)
}

func setRouters(e *echo.Echo, s *service.Service, tasks chan *service.Item, logger *zap.Logger) {

	e.POST("/api/v1/tasks/", handlers.StartProcessing(s, tasks, logger))

}
