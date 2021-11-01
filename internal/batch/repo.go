package batch

import "context"

type Repo interface {
	GetAllBatches(ctx context.Context) ([]Batch, error)
	GetBatch(ctx context.Context, id string) (*Batch, error)
}
