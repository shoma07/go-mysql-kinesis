package producer

import (
	"log"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type EventHandler struct {
	events    *chan interface{}
	positions *chan mysql.Position
}

func NewEventHandler(events *chan interface{}, positions *chan mysql.Position) (*EventHandler, error) {
	h := &EventHandler{}
	h.events = events
	h.positions = positions
	return h, nil
}

func (h *EventHandler) OnRotate(*replication.RotateEvent) error {
	return nil
}

func (h *EventHandler) OnTableChanged(schema, table string) error {
	log.Printf("OnTableChanged: schema=%s, table=%s", schema, table)
	return nil
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	log.Printf("OnRow: Table=%s Action=%s Rows=%s Header=%s\n", e.Table, e.Action, e.Rows, e.Header)
	*h.events <- e
	return nil
}

func (h *EventHandler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	log.Printf("OnPosSynced: %s\n", pos)
	*h.positions <- pos

	return nil
}

func (h *EventHandler) OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

func (h *EventHandler) OnXID(mysql.Position) error {
	return nil
}

func (h *EventHandler) OnGTID(mysql.GTIDSet) error {
	return nil
}

func (h *EventHandler) String() string {
	return "EventHandler"
}
