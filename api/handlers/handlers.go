package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"kznexp/service"
)

type Service interface {
	Process(ctx context.Context, batch *service.Batch, tasks chan *service.Item) error
}

func StartProcessing(s Service, tasks chan *service.Item, logger *zap.Logger) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {

		type reqTasks struct {
			CountTasks int `json:"count_tasks"`
		}

		req := new(reqTasks)
		if err := ctx.Bind(&req); err != nil {
			return responseError(ctx, logger, http.StatusBadRequest, err)
		}

		if req == nil || req.CountTasks == 0 {
			return responseError(ctx, logger, http.StatusBadRequest, errors.New("nothing to process"))
		}

		batch := batch(req.CountTasks)

		if err := s.Process(ctx.Request().Context(), &batch, tasks); err != nil {
			return responseError(ctx, logger, http.StatusInternalServerError, err)
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func batch(count int) service.Batch {
	b := make([]*service.Item, 0)
	for i := 0; i < count; i++ {
		b = append(b, &service.Item{})
	}
	return b
}

func responseError(ctx echo.Context, logger *zap.Logger, code int, err error) error {
	logger.Error("error: ", zap.Error(err))
	return ctx.String(code, err.Error())
}
