package actions

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/actions/handlers"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/types"
)

// TODO: This should not be an action. Maybe introduce Job type?
type HttpServer struct {
	Addr string
	Port int
}

type HandlerContext struct {
	Context
}

func (c HandlerContext) ReadEvents(stream string, cursor string, amount int) ([]types.Event, error) {
	re := ReadEvents{Stream: stream, Cursor: cursor, Amount: amount}
	err := re.Run(c)
	return re.Events, errors.Wrap(err, "HandlerContext ReadEvents .Run")
}

func (c HandlerContext) ListStreams() ([]string, error) {
	session, err := c.Persistence().Session()
	if err != nil {
		return []string{}, errors.Wrap(err, "HandlerContext ListStreams .Persistence.Session")
	}
	return session.ListStreams()
}

func (hs HttpServer) Run(c Context) error {
	handler := c.HttpHandler()
	jobs := c.BackgroundJobs()
	handler_context := handlers.Handler{HandlerContext{c}}
	handler.Get("/_check", func(ctx http.GetContext) { ctx.SendJson("OK") })
	handler.Get("/streams", handler_context.GetStreams)
	handler.Get("/stream/:id", handler_context.GetStream)
	handler.Post("/bulk/stream/:id", func(ctx http.PostContext) {
		var err error
		stream := ctx.GetSegment("id")

		job := BulkWriteJob{}
		err = ctx.ReadJSON(&job)

		if err != nil {
			ctx.ServerErr(fmt.Errorf("Expected body, encountered: %v", err))
			return
		}

		job.Stream = stream

		jobs <- job
		c.Telemetry().Gauge("jobs.len", []string{}, len(jobs))
		ctx.SendJson("OK")
	})
	handler.Post("/stream/:id", func(ctx http.PostContext) {
		var err error
		stream := ctx.GetSegment("id")
		body, err := ctx.JsonBodyParams()
		if err != nil {
			ctx.ServerErr(fmt.Errorf("Expected body, encountered: %v", err))
			return
		}

		jobs <- WriteJob{Stream: stream, Payload: body}
		c.Telemetry().Gauge("jobs.len", []string{}, len(jobs))
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
