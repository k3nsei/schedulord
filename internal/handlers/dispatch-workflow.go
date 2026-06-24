package handlers

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	gh "github.com/k3nsei/schedulord/internal/github"
	"github.com/k3nsei/schedulord/internal/models"
)

type dispatchWorkflowHandler struct {
	models.TaskHandler
}

func (h dispatchWorkflowHandler) Run(ctx context.Context, payload any) error {
	token, ok := gh.GetGithubToken(ctx)
	if !ok {
		return fmt.Errorf("GitHub token is missing")
	} else if token == "" {
		return fmt.Errorf("GitHub token is empty")
	}

	var p *models.DispatchWorkflowPayload
	err := mapstructure.Decode(payload, &p)
	if err != nil {
		return fmt.Errorf("failed to decode payload")
	}

	err = gh.DispatchWorkflow(ctx, token, p.Owner, p.Repository, p.Ref, p.WorkflowID, p.Inputs)
	if err != nil {
		return err
	}

	return nil
}
