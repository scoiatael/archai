package http

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func NewIris(c Context, useDevLogger bool) *IrisHandler {
	handler := IrisHandler{}
	handler.context = c

	app := iris.New()
	if useDevLogger {
		app.Adapt(
			// adapt a logger which prints all errors to the os.Stdout
			iris.DevLogger(),
		)
	}
	app.Adapt(
		// adapt the adaptors/httprouter or adaptors/gorillamux
		httprouter.New(),
	)

	handler.framework = app
	return &handler
}

func NewFastHttp(c Context) *FastHttpHandler {
	return &FastHttpHandler{Context: c}
}
