package actions

import (
	"github.com/scoiatael/archai/persistence"
)

type Context interface {
	Persistence() persistence.Provider
	Migrations() map[string]persistence.Migration
	Version() string
}

type Action interface {
	Run(Context) error
}
