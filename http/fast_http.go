package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type FastHttpContext struct {
	*fasthttp.RequestCtx
	Context Context
}

func (hc FastHttpContext) ServerErr(err error) {
	hc.Context.HandleErr(err)
	hc.Error(fmt.Sprintf(`{ "error": "%v" }`, err), fasthttp.StatusInternalServerError)
}

func (hc FastHttpContext) SendJson(response interface{}) {
	dump, err := json.Marshal(response)
	if err != nil {
		hc.ServerErr(err)
	} else {
		hc.SetBody(dump)
	}
}

// TODO: Do this normal way
func (hc FastHttpContext) GetSegment(index string) string {
	segments := strings.Split(string(hc.Path()), "/")
	if len(segments) == 0 {
		return ""
	} else {
		return segments[len(segments)-1]
	}
}

type FastHttpGetContext struct {
	FastHttpContext
}

func (gc FastHttpGetContext) StringParam(name string) string {
	val := gc.QueryArgs().Peek(name)
	return string(val)
}

func (gc FastHttpGetContext) IntParam(name string, def int) int {
	val := gc.QueryArgs().Peek(name)
	i, err := strconv.Atoi(string(val))
	if err != nil {
		return def
	}
	return i
}

type FastHttpPostContext struct {
	FastHttpContext
}

const expectedJSON = `{ "error": "expected JSON body" }`

func (pc FastHttpPostContext) JsonBodyParams() (map[string]interface{}, error) {
	body := pc.PostBody()
	read := make(map[string]interface{})
	err := json.Unmarshal(body, &read)
	if err != nil {
		pc.Error(expectedJSON, fasthttp.StatusBadRequest)
	}
	return read, err
}

// TODO: Add routing ;)
type FastHttpHandlers struct {
	POST func(PostContext)
	GET  func(GetContext)
}

type FastHttpHandler struct {
	handlers FastHttpHandlers
	Context  Context
}

func (h *FastHttpHandler) Get(path string, handler func(GetContext)) {
	h.handlers.GET = handler
}

func (h *FastHttpHandler) Post(path string, handler func(PostContext)) {
	h.handlers.POST = handler
}

const (
	methodNotAllowed = `{ "error": "method not allowed" }`
	contentType      = `application/json`
)

func (h *FastHttpHandler) compile() func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType(contentType)
		httpCtx := FastHttpContext{ctx, h.Context}
		if ctx.IsPost() {
			h.handlers.POST(FastHttpPostContext{httpCtx})
		} else if ctx.IsGet() {
			h.handlers.GET(FastHttpGetContext{httpCtx})
		} else {
			ctx.Error(methodNotAllowed, fasthttp.StatusMethodNotAllowed)
		}
	}
}

func (h *FastHttpHandler) Run(addr string) error {
	return errors.Wrap(fasthttp.ListenAndServe(addr, h.compile()), "Starting fasthttp server")
}
