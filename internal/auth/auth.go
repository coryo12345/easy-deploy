package auth

type AuthRepository interface {
	Authenticate(password string) bool
}

type auth struct {
	password string
}

func New(password string) AuthRepository {
	return &auth{
		password: password,
	}
}

func (a auth) Authenticate(password string) bool {
	return password == a.password
}
