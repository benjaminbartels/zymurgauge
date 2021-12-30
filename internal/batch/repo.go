package batch

import "context"

type Repo interface {
	GetAll(ctx context.Context) ([]Batch, error)
	Get(ctx context.Context, id string) (*Batch, error)
}
