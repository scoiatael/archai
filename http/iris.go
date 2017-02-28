package http

import (
	"gopkg.in/kataras/iris.v6"
)

type IrisHttpContext struct {
	*iris.Context
	context Context
}

func (hc IrisHttpContext) SendJson(response interface{}) {
	hc.JSON(iris.StatusOK, response)
}

func (hc IrisHttpContext) ServerErr(err error) {
	hc.context.HandleErr(err)
	hc.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
}

func (hc IrisHttpContext) GetSegment(index string) string {
	return hc.Param(index)
}

type IrisGetContext struct {
	IrisHttpContext
}

func (hc IrisHttpContext) StringParam(index string) string {
	return hc.Param(index)
}

func (hc IrisHttpContext) IntParam(index string, def int) int {
	val, err := hc.ParamInt(index)
	if err != nil {
		return def
	}

	return val
}

type IrisPostContext struct {
	IrisHttpContext
}

func (hc IrisPostContext) JsonBodyParams() (map[string]interface{}, error) {
	sess := iris.Map{}
	err := hc.ReadJSON(&sess)

	if err != nil {
		hc.JSON(iris.StatusBadRequest,
			iris.Map{"error": "expected JSON body",
				"details":  err,
				"received": hc.Request.Body,
			})
	}
	return sess, err
}

type IrisHandler struct {
	framework *iris.Framework
	context   Context
}

func (h *IrisHandler) Get(path string, handler func(GetContext)) {
	h.framework.Get(path, func(ctx *iris.Context) {
		handler(IrisGetContext{IrisHttpContext{ctx, h.context}})
	})
}

func (h *IrisHandler) Post(path string, handler func(PostContext)) {
	h.framework.Post(path, func(ctx *iris.Context) {
		handler(IrisPostContext{IrisHttpContext{ctx, h.context}})
	})
}

func (h *IrisHandler) Run(addr string) error {
	h.framework.Listen(addr)

	return nil
}
