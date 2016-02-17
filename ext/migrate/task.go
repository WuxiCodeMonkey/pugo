package migrate

import (
	"bytes"
	"github.com/go-xiaohei/pugo/app/builder"
)

var (
	manager *Manager
)

func init() {
	manager = NewManager(new(RSS))
}

type (
	// Task define migration methods
	Task interface {
		Name() string
		Detect(*builder.Context) (Task, error)
		Action(*builder.Context) (map[string]*bytes.Buffer, error)
	}
	// Manager manage tasks in global
	Manager struct {
		tasks map[string]Task
	}
)

// NewManager new manager with typed tasks
func NewManager(tasks ...Task) *Manager {
	m := &Manager{
		tasks: make(map[string]Task),
	}
	for _, t := range tasks {
		m.tasks[t.Name()] = t
	}
	return m
}