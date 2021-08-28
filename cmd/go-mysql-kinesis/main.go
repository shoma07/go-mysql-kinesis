package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type MyPosition struct {
	Name string `json: "name,string"`
	Pos  uint32 `json: "pos"`
}

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Printf("OnRow: Table=%s Action=%s Rows=%s Header=%s\n", e.Table, e.Action, e.Rows, e.Header)
	return nil
}

func (h *MyEventHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	log.Printf("OnPosSynced: %s\n", pos)
	myPosition := MyPosition{pos.Name, pos.Pos}
	bytes, err := json.Marshal(&myPosition)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile("position.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, string(bytes))

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
		log.Fatal(err)
	}

	c.SetEventHandler(&MyEventHandler{})

	var initPosition MyPosition
	bytes, err := ioutil.ReadFile("position.json")
	if err != nil {
		log.Printf("Warning init position: %s\n", err)
	} else {
		if loadErr := json.Unmarshal(bytes, &initPosition); err != nil {
			panic(loadErr)
		}
	}

	if initPosition.Name == "" && initPosition.Pos == 0 {
		c.Run()
	} else {
		log.Printf("init position: %s\n", initPosition)
		c.RunFrom(mysql.Position{initPosition.Name, initPosition.Pos})
	}
}
