package jwtUtil

type Login struct {
	ID       uint64
	Account  string
	Password string
	Roles    string
	Enable   bool
	IP       string
}

type Token struct{}

func (Token) LoginToken(user *Login) (token string, claims CustomClaims, err error) {
	jwtManager := NewJWT()
	claims = jwtManager.CreateClaims(BaseClaims{
		ID:       user.ID,
		Account:  user.Account,
		Password: user.Password,
		Roles:    user.Roles,
		Enable:   user.Enable,
		IP:       user.IP,
	})
	token, err = jwtManager.CreateToken(claims)
	return
}
