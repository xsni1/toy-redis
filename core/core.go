package core

import (
	"github.com/xsni1/toy-redis/parser"
	"github.com/xsni1/toy-redis/store"
)

type Core struct {
	store *store.Store
}

func NewCore(store *store.Store) *Core {
	return &Core{
		store: store,
	}
}

func (c *Core) Execute(cmd parser.ParsedMessage) {}
