package handlers

import (
	"github.com/pkg/errors"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/simplejson"
)

func (gs Handler) GetStreams(ctx http.GetContext) {
	streams, err := gs.Context.ListStreams()
	if err != nil {
		ctx.ServerErr(errors.Wrap(err, "GetStreams Handle .ListStreams"))
		return
	}
	view := make(simplejson.Object)
	view["streams"] = streams
	ctx.SendJson(view)
}
