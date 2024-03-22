package storage

import "context"
import "slackwatch/backend/internal/model"

type Storage interface {
	Save(ctx context.Context, data model.Data) error
	Get(ctx context.Context, id string) (model.Data, error)
}
