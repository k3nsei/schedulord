package models

type TaskType string

const (
	DispatchWorkflow TaskType = "dispatch-workflow"
)

func (tt TaskType) String() string {
	return string(tt)
}
