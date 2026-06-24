package models

import "context"

type TaskHandler interface {
	Run(ctx context.Context, payload any) error
}
