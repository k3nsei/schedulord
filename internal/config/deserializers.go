package config

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/k3nsei/schedulord/internal/models"
)

var registry = map[string]reflect.Type{
	string(models.DispatchWorkflow): reflect.TypeOf(models.DispatchWorkflowTask{}),
}

func decodeTask(raw map[string]any) (models.Task, error) {
	tt, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing task type")
	}

	typ, ok := registry[tt]
	if !ok {
		return nil, fmt.Errorf("unknown task type %q", tt)
	}

	v := reflect.New(typ)

	if err := mapstructure.Decode(raw, v.Interface()); err != nil {
		return nil, err
	}

	t, ok := v.Interface().(models.Task)
	if !ok {
		return nil, fmt.Errorf("%s does not implement Task", typ.Name())
	}

	return t, nil
}

type TaskSpec struct {
	models.Task
}

func (ts *TaskSpec) UnmarshalMapstructure(input any) error {
	raw, ok := input.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map, got %T", input)
	}

	t, err := decodeTask(raw)
	if err != nil {
		return err
	}

	ts.Task = t

	return nil
}
