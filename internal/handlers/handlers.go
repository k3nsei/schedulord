package handlers

import (
	"github.com/k3nsei/schedulord/internal/models"
)

var handlers = map[models.TaskType]models.TaskHandler{
	models.DispatchWorkflow: dispatchWorkflowHandler{},
}

func ResolveHandler(tt models.TaskType) (models.TaskHandler, bool) {
	handler, ok := handlers[tt]
	if !ok {
		return nil, false
	}
	return handler, true
}
