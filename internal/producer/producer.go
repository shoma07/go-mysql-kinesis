package producer

import (
	"log"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

// NOTE:
type Producer struct {
	canal     *canal.Canal
	pos       *Position
	handler   *EventHandler
	events    chan interface{}
	positions chan mysql.Position
}

func NewProducer() (*Producer, error) {
	p := &Producer{}
	p.events = make(chan interface{}, 4096)
	p.positions = make(chan mysql.Position, 4096)

	var err error

	p.pos, err = NewPosition()
	if err != nil {
		log.Printf("[ERROR] Catch Error in NewProducer: %s\n", err)
		return nil, err
	}

	p.handler, err = NewEventHandler(&p.events, &p.positions)
	if err != nil {
		log.Printf("[ERROR] Catch Error in NewProducer: %s\n", err)
		return nil, err
	}

	cfg := canal.NewDefaultConfig()
	cfg.Addr = "db:3306"
	cfg.User = "root"
	cfg.Password = "root"
	cfg.Dump.TableDB = "chat"
	cfg.Dump.Tables = []string{"messages"}

	p.canal, err = canal.NewCanal(cfg)
	if err != nil {
		log.Printf("[ERROR] Initialize Canal: %s", err)
		return nil, err
	}
	p.canal.SetEventHandler(p.handler)

	return p, nil
}

func (p *Producer) Run() error {
	if p.pos.IsNew() {
		p.canal.Run()
	} else {
		p.canal.RunFrom(p.pos.MySQLPosition())
	}
	return nil
}
