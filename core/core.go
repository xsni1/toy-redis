package core

import (
	"fmt"

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

func (c *Core) Execute(cmd parser.ParsedMessage) {
	if cmd.Msgtype != parser.Array {
		fmt.Printf("wrong type\n")
		return
	}

    fmt.Print("ARGS, ", cmd.Args, "\n")
}
