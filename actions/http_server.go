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

type WriteJob struct {
	payload simplejson.Object
	stream  string
}

func writer(jobs <-chan WriteJob, c Context) {
	for j := range jobs {
		j.Run(c)
	}
}

func (wj *WriteJob) Run(c Context) {
	payload, err := json.Marshal(wj.payload)
	err = errors.Wrap(err, "HTTP server marshalling payload to write event")
	if err != nil {
		c.HandleErr(err)
		return
	}
	action := WriteEvent{Stream: wj.stream, Payload: payload, Meta: make(map[string]string)}
	action.Meta["origin"] = "http"
	action.Meta["compressed"] = "false"
	err = action.Run(c)
	if err != nil {
		c.HandleErr(errors.Wrap(err, "HTTP server writing event"))
	} else {
		c.Telemetry().Incr("write", []string{"stream:" + wj.stream})
	}
}

func (hs HttpServer) Run(c Context) error {
	handler := c.HttpHandler()
	jobs := make(chan WriteJob, 50)
	for w := 0; w < c.Concurrency(); w++ {
		go writer(jobs, c)
	}
	handler.Get("/_check", func(ctx http.GetContext) {
		ctx.SendJson("OK")
	})
	handler.Get("/stream/:id", func(ctx http.GetContext) {
		stream := ctx.GetSegment("id")
		action := ReadEvents{Stream: stream}
		action.Cursor = ctx.StringParam("cursor")
		action.Amount = ctx.IntParam("amount", 10)
		err := errors.Wrap(action.Run(c), "HTTP server reading events")
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
				err = errors.Wrap(err, "HTTP server marshalling response with read events")
				if err != nil {
					c.HandleErr(err)
					ctx.ServerErr(err)
				}
				events[i]["blob"] = payload
			}
			root["results"] = events
			c.Telemetry().Incr("read", []string{"stream:" + stream})
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

		jobs <- WriteJob{stream: stream, payload: body}
		ctx.SendJson("OK")
	})

	connString := fmt.Sprintf("%s:%d", hs.Addr, hs.Port)
	return errors.Wrap(handler.Run(connString), "HttpServer starting..")
}

func (hs HttpServer) Stop() {
}

func (hs HttpServer) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("Start HTTP server on %s:%d", hs.Addr, hs.Port))
}
