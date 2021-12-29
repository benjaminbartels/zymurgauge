package chamber

import "context"

type Controller interface {
	Repo
	StartFermentation(ctx context.Context, chamberID string, step int) error
	StopFermentation(chamberID string) error
}

type Repo interface {
	GetAll() ([]*Chamber, error) // TODO: add ctx? // TODO: should this return a slice of pointers?
	Get(id string) (*Chamber, error)
	Save(c *Chamber) error
	Delete(id string) error
}
