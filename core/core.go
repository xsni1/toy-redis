package core

import (
	"fmt"
	"strings"

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

func (c *Core) Execute(cmd parser.ParsedMessage) []byte {
	if cmd.Msgtype != parser.Array {
		return []byte{}
	}

	if len(cmd.Args) == 0 {
		fmt.Printf("empty\n")
		return []byte{}
	}

	switch strings.ToLower(cmd.Args[0]) {
	case "set":
		c.store.M.Store(cmd.Args[1], cmd.Args[2])
		return parser.SimpleStringReply("OK")
	case "get":
		val, ok := c.store.M.Load(cmd.Args[1])
		if !ok {
			return parser.NullReply()
		}
		return parser.SimpleStringReply(val.(string))
	}

	return parser.NullReply()
}
