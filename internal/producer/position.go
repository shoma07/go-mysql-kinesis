package producer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-mysql-org/go-mysql/mysql"
)

type Position struct {
	Name string `json: "name,string"`
	Pos  uint32 `json: "pos"`
}

func NewPosition() (*Position, error) {
	pos := &Position{}
	bytes, err := ioutil.ReadFile("position.json")
	if err != nil {
		log.Printf("[WARN] Initialize Position Error: %s\n", err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, pos); err != nil {
		log.Printf("[ERROR] Initialize Position Error: %s\n", err)
		return nil, err
	}

	return pos, nil
}

func (pos *Position) IsNew() bool {
	return pos.Name == "" && pos.Pos == 0
}

func (pos *Position) MySQLPosition() mysql.Position {
	return mysql.Position{pos.Name, pos.Pos}
}

func (pos *Position) Save(mp mysql.Position) error {
	pos.Name = mp.Name
	pos.Pos = mp.Pos

	bytes, err := json.Marshal(&pos)
	if err != nil {
		log.Printf("[ERROR] Position Dump Error: %s\n", err)
		return err
	}

	file, err := os.OpenFile("position.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("[ERROR] File Open Error: %s\n", err)
		return err
	}

	defer file.Close()
	fmt.Fprintln(file, string(bytes))

	return nil
}
