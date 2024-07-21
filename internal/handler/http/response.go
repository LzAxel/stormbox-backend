package http

import "github.com/labstack/echo/v4"

type DataHTTPResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

type ErrorHTTPResponse struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error"`
}

func ProcessDataHTTPResponse(ctx echo.Context, status int, data any) error {
	return ctx.JSON(status, DataHTTPResponse{
		Status: status,
		Data:   data,
	})
}
