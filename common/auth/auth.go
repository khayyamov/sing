package auth

import "github.com/sagernet/sing/common"

type User struct {
	Username string
	Password string
}

type Authenticator struct {
	UserMap map[string][]string
}

func NewAuthenticator(users []User) *Authenticator {
	if len(users) == 0 {
		return nil
	}
	au := &Authenticator{
		UserMap: make(map[string][]string),
	}
	for _, user := range users {
		au.UserMap[user.Username] = append(au.UserMap[user.Username], user.Password)
	}
	return au
}

func (au *Authenticator) AddUserToAuthenticator(users []User) {
	if len(users) > 0 {
		for _, user := range users {
			au.UserMap[user.Username] = append(au.UserMap[user.Username], user.Password)
		}
	}
}

func (au *Authenticator) DeleteUserToAuthenticator(users []User) {
	if len(users) > 0 {
		for _, user := range users {
			delete(au.UserMap, user.Username)
		}
	}
}

func (au *Authenticator) Verify(username string, password string) bool {
	passwordList, ok := au.UserMap[username]
	return ok && common.Contains(passwordList, password)
}
