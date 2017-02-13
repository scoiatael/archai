package actions

import (
	"github.com/scoiatael/archai/persistence"
)

type Context interface {
	Persistence() persistence.Provider
	Migrations() map[string]persistence.Migration
}

type Action interface {
	Run(Context) error
}
