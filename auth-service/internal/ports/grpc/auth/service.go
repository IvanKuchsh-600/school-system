package auth

type AuthService interface {
	Register(email, password, role string) (string, error)
	Login(email, password string) (string, error)
}
