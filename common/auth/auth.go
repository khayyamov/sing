package auth

import "github.com/sagernet/sing/common"

type User struct {
	Username string
	Password string
}

type Authenticator struct {
	userMap map[string][]string
}

func NewAuthenticator(users []User) *Authenticator {
	if len(users) == 0 {
		return nil
	}
	au := &Authenticator{
		userMap: make(map[string][]string),
	}
	for _, user := range users {
		au.userMap[user.Username] = append(au.userMap[user.Username], user.Password)
	}
	return au
}

func (au *Authenticator) AddUserToAuthenticator(users []User) {
	if len(users) > 0 {
		for _, user := range users {
			au.userMap[user.Username] = append(au.userMap[user.Username], user.Password)
		}
	}
}

func (au *Authenticator) DeleteUserToAuthenticator(users []User) {
	if len(users) > 0 {
		for _, user := range users {
			delete(au.userMap, user.Username)
		}
	}
}

func (au *Authenticator) Verify(username string, password string) bool {
	passwordList, ok := au.userMap[username]
	return ok && common.Contains(passwordList, password)
}
