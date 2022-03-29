package auth

import "time"

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	ModTime  time.Time `json:"modTime"`
}
