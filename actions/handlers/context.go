package handlers

import (
	"github.com/scoiatael/archai/telemetry"
	"github.com/scoiatael/archai/types"
)

type Context interface {
	ReadEvents(string, string, int) ([]types.Event, error)
	ListStreams() ([]string, error)
	Telemetry() telemetry.Datadog
}

type Handler struct {
	Context
}
