package actions

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/simplejson"
)

type HttpServer struct {
	Addr string
	Port int
}

func (hs HttpServer) Run(c Context) error {
	handler := c.HttpHandler()
	handler.Get("/stream/:id", func(ctx http.GetContext) {
		stream := ctx.GetSegment("id")
		action := ReadEvents{Stream: stream}
		action.Cursor = ctx.StringParam("cursor")
		action.Amount = ctx.IntParam("amount", 10)
		err := action.Run(c)
		if err != nil {
			c.HandleErr(err)
			ctx.ServerErr(err)
		} else {
			root := make(simplejson.Object)
			events := make(simplejson.ObjectArray, len(action.Events))
			for i, ev := range action.Events {
				events[i] = make(simplejson.Object)
				events[i]["ID"] = ev.ID
				payload, err := simplejson.Read(ev.Blob)
				if err != nil {
					c.HandleErr(err)
					ctx.ServerErr(err)
				}
				events[i]["blob"] = payload
			}
			root["results"] = events
			ctx.SendJson(root)
		}
	})
	handler.Post("/stream/:id", func(ctx http.PostContext) {
		var err error
		stream := ctx.GetSegment("id")
		body, err := ctx.JsonBodyParams()
		if err != nil {
			// Error was already sent
			return
		}
		payload, err := json.Marshal(body)
		if err != nil {
			c.HandleErr(err)
			ctx.ServerErr(err)
			return
		}
		action := WriteEvent{Stream: stream, Payload: payload, Meta: make(map[string]string)}
		action.Meta["origin"] = "http"
		action.Meta["compressed"] = "false"
		err = action.Run(c)
		if err != nil {
			c.HandleErr(err)
			ctx.ServerErr(err)
		} else {
			ctx.SendJson("OK")
		}
	})

	connString := fmt.Sprintf("%s:%d", hs.Addr, hs.Port)
	return errors.Wrap(handler.Run(connString), "HttpServer starting..")
}

func (hs HttpServer) Stop() {
}
