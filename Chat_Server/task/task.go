package task

import (
	"github.com/google/uuid"
	"github.com/phuslu/log"
)

const (
	UNICAST   = "unicast"
	MULTICAST = "multicast"
	BROADCAST = "broadcast"
)

type Task struct {
	Id     uuid.UUID
	Type   string
	Weight float64
	Task   string
}

func (t *Task) MarshalObject(e *log.Entry) {
	e.Str("id", t.Id.String()).Str("type", t.Type).Float64("weight", t.Weight).Str("task", t.Task)
}

func NewTask(_type, task string, weight float64) *Task {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &Task{
		Id:     uid,
		Type:   _type,
		Weight: weight, // TODO - determine weight of task somehow, may be useful
		Task:   task,
	}
}
