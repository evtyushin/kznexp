package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"kznexp/service"
)

func TestStartProcessing(t *testing.T) {
	type reqBody struct {
		CountTasks int `json:"count_tasks"`
	}

	tests := []struct {
		name     string
		body     reqBody
		prepare  func(s *MockService)
		wantCode int
		wantBody string
	}{
		{
			name: "200Code",
			body: reqBody{CountTasks: 10},
			prepare: func(s *MockService) {
				s.On("Process", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			wantCode: 200,
			wantBody: "",
		},
		{
			name:     "400CodeEmptyBody",
			wantCode: 400,
			wantBody: "nothing to process",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := new(MockService)
			if tt.prepare != nil {
				tt.prepare(s)
			}
			e := echo.New()

			bs, _ := json.Marshal(&tt.body)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bs))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			l := zap.NewNop()
			ch := make(chan *service.Item)

			_ = StartProcessing(s, ch, l)(c)

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, rec.Body.String())
			mock.AssertExpectationsForObjects(t, s)
		})

	}
}
