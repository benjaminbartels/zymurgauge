package chamber

type Repo interface {
	GetAllChambers() ([]Chamber, error) // TODO: add ctx?
	GetChamber(id string) (*Chamber, error)
	SaveChamber(c *Chamber) error
	DeleteChamber(id string) error
}
