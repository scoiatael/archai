package handlers

import (
	"github.com/pkg/errors"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/simplejson"
	"github.com/scoiatael/archai/types"
)

func serializeEvents(events []types.Event) (simplejson.Object, error) {
	root := make(simplejson.Object)
	results := make([]simplejson.Object, len(events))
	cursor := make(simplejson.Object)
	for i, ev := range events {
		payload, err := simplejson.Read(ev.Blob)
		if err != nil {
			return root, errors.Wrap(err, "HTTP server marshalling response with read events")
		}
		results[i] = payload
		cursor["next"] = ev.ID
	}
	root["results"] = results
	root["cursor"] = cursor
	return root, nil
}

func (gs Handler) GetStream(ctx http.GetContext) {
	stream := ctx.GetSegment("id")
	events, err := gs.Context.ReadEvents(
		stream,
		ctx.StringParam("cursor"),
		ctx.IntParam("amount", 10),
	)
	if err != nil {
		ctx.ServerErr(errors.Wrap(err, "GetStream Handle ReadEvents"))
		return
	}
	json, err := serializeEvents(events)
	if err != nil {
		ctx.ServerErr(err)
	}
	gs.Context.Telemetry().Incr("read", []string{"stream:" + stream})
	ctx.SendJson(json)
}
