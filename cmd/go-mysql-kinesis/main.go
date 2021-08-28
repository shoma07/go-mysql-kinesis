package main

import (
	"log"

	"github.com/go-mysql-org/go-mysql/canal"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Printf("%s %s\n", e.Action, e.Rows)
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func main() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "db:3306"
	cfg.User = "root"
	cfg.Password = "root"
	cfg.Dump.TableDB = "chat"
	cfg.Dump.Tables = []string{"messages"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Printf("%s", err)
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&MyEventHandler{})

	// Start canal
	c.Run()
}
