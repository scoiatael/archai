package actions

import (
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/persistence"
)

type HttpHandler interface {
	Get(string, func(http.GetContext))
	Post(string, func(http.PostContext))
	Run(string) error
}

type Context interface {
	Persistence() persistence.Provider
	Migrations() map[string]persistence.Migration
	Version() string
	HandleErr(error)
	HttpHandler() HttpHandler
}

type Action interface {
	Run(Context) error
}
