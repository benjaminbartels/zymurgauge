package internal

import "time"

// Device is a heating or cooling element that can be switched on or off via a GPIO pin.
type Device struct {
	GPIO     int           `json:"gpio"`
	Cooldown time.Duration `json:"cooldown"`
}

// Equals returns true if Device d is equal to the given Device o
// ToDo: Re-assess equality (https://golangbot.com/structs/?utm_source=golangweekly&utm_medium=email)
func (d *Device) Equals(o *Device) bool {
	if d.GPIO != o.GPIO {
		return false
	} else if d.Cooldown != o.Cooldown {
		return false
	}

	return true
}
