package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func DispatchWorkflow(ctx context.Context, token, owner, repository, ref, workflowId string, inputs any) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches",
		owner, repository, workflowId,
	)

	if ref == "" {
		ref = "main"
	}

	body, _ := json.Marshal(map[string]any{
		"ref":    ref,
		"inputs": inputs,
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("github error: %s", resp.Status)
	}

	return nil
}
