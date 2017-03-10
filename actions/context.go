package actions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/persistence"
	"github.com/scoiatael/archai/telemetry"
)

type HttpHandler interface {
	Get(string, func(http.GetContext))
	Post(string, func(http.PostContext))
	Run(string) error
	Stop(context.Context)
}

type Context interface {
	Persistence() persistence.Provider
	Migrations() map[string]persistence.Migration
	Version() string
	HandleErr(error)
	HttpHandler() HttpHandler
	Telemetry() telemetry.Datadog
	Concurrency() int
	Retries() int
	Backoff(int) time.Duration
}

type Action interface {
	Run(Context) error
	json.Marshaler
}
