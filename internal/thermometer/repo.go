package thermometer

type Repo interface {
	GetThermometerIDs() ([]string, error)
}
