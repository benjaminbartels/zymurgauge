package safeclose

import "io"

// Close closes the io.Closer and sets err if an error occurs while closing
func Close(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
