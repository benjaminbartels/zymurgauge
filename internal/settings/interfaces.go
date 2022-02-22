package settings

type Repo interface {
	Get() (*Settings, error)
	Save(c *Settings) error
}
