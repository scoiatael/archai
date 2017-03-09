package actions

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/simplejson"
)

// TODO: This should not be an action. Maybe introduce Job type?
type HttpServer struct {
	Addr string
	Port int
}

type BackgroundJob interface {
	Run(c Context)
}

func persistEvent(stream string, payload []byte, origin string, c Context) error {
	action := WriteEvent{Stream: stream, Payload: payload, Meta: make(map[string]string)}
	action.Meta["origin"] = origin
	action.Meta["compressed"] = "false"
	action.Meta["time"] = string(time.Now().Unix())
	return action.Run(c)
}

type WriteJob struct {
	payload simplejson.Object
	stream  string
}

func (wj WriteJob) Run(c Context) {
	payload, err := json.Marshal(wj.payload)
	err = errors.Wrap(err, "HTTP server marshalling payload to write event")
	if err != nil {
		c.HandleErr(err)
		return
	}
	err = persistEvent(wj.stream, payload, "http; write_job", c)
	if err != nil {
		c.HandleErr(errors.Wrap(err, "HTTP server writing event"))
	} else {
		c.Telemetry().Incr("write", []string{"stream:" + wj.stream})
	}
}

type BulkWriteJob struct {
	schema  []interface{}
	objects []interface{}
	stream  string
}

func makeObjectWithSchema(obj interface{}, schema []interface{}) (simplejson.Object, error) {
	object_with_schema := make(simplejson.Object)
	object, conv := obj.([]interface{})
	if !conv {
		return object_with_schema, fmt.Errorf("Failed to convert obj to array")
	}
	for j, name := range schema {
		name, conv := name.(string)
		if !conv {
			return object_with_schema, fmt.Errorf("%d: Failed to convert schema value to string", j)
		}
		if len(object) <= j {
			return object_with_schema, fmt.Errorf("%d: Not enough values", j)
		}
		object_with_schema[name] = object[j]
	}
	return object_with_schema, nil
}

func (wj BulkWriteJob) Run(c Context) {
	c.Telemetry().Incr("bulk_write.aggregate", []string{"stream:" + wj.stream})
	for i, obj := range wj.objects {
		object, err := makeObjectWithSchema(obj, wj.schema)
		err = errors.Wrap(err, fmt.Sprintf("HTTP server splitting payload to bulk_write event at %d", i))
		if err != nil {
			c.HandleErr(err)
			return
		}
		payload, err := json.Marshal(object)
		err = errors.Wrap(err, fmt.Sprintf("HTTP server marshalling payload to bulk_write event at %d", i))
		if err != nil {
			c.HandleErr(err)
			return
		}
		err = persistEvent(wj.stream, payload, "http; bulk_write_job", c)
		if err != nil {
			c.HandleErr(errors.Wrap(err, "HTTP server bulk_writing events"))
		} else {
			c.Telemetry().Incr("bulk_write.singular", []string{"stream:" + wj.stream})
		}
	}
}

func writer(jobs <-chan BackgroundJob, c Context) {
	for j := range jobs {
		j.Run(c)
	}
}

func (hs HttpServer) Run(c Context) error {
	handler := c.HttpHandler()
	jobs := make(chan BackgroundJob, 50)
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
	handler.Post("/bulk/stream/:id", func(ctx http.PostContext) {
		var err error
		stream := ctx.GetSegment("id")
		body, err := ctx.JsonBodyParams()
		if err != nil {
			// Error was already sent
			return
		}
		objects, conv := body["data"].([]interface{})
		if !conv {
			ctx.ServerErr(fmt.Errorf("'data' field is not an Array (is %v)", reflect.TypeOf(body["data"])))
			return
		}
		schema, conv := body["schema"].([]interface{})
		if !conv {
			ctx.ServerErr(fmt.Errorf("'schema' field is not an Array (is %v)", reflect.TypeOf(body["schema"])))
			return
		}

		jobs <- BulkWriteJob{stream: stream, objects: objects, schema: schema}
		ctx.SendJson("OK")
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
