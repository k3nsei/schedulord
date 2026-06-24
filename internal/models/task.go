package models

type BaseTask struct {
	Type     TaskType `mapstructure:"type"`
	Cron     string   `mapstructure:"cron"`
	Disabled bool     `mapstructure:"disabled"`
	Payload  any      `mapstructure:"payload"`
}

type Task interface {
	GetBase() BaseTask
	isTask()
}

type DispatchWorkflowTask struct {
	BaseTask `mapstructure:",squash"`
	Payload  DispatchWorkflowPayload `mapstructure:"payload"`
}

type DispatchWorkflowPayload struct {
	Owner      string            `mapstructure:"owner"`
	Repository string            `mapstructure:"repo"`
	Ref        string            `mapstructure:"ref"`
	WorkflowID string            `mapstructure:"workflow"`
	Inputs     map[string]string `mapstructure:"inputs"`
}

func (t DispatchWorkflowTask) GetBase() BaseTask {
	return t.BaseTask
}

func (t DispatchWorkflowTask) isTask() {}
