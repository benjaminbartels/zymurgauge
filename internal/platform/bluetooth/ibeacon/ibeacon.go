package ibeacon

type EventType uint8

const (
	Online EventType = iota
	Offline
)

func (t EventType) String() string {
	switch t {
	case Online:
		return "online"
	case Offline:
		return "offline"
	}

	return "unknown"
}

type Event struct {
	UUID    string
	Type    EventType
	IBeacon IBeacon
}

type IBeacon struct {
	ProximityUUID string
	Major         uint16
	Minor         uint16
	MeasuredPower uint16
}
