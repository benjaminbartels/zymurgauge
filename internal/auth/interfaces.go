package auth

type UserRepo interface {
	Get() (*User, error)
	Save(u *User) error
}
