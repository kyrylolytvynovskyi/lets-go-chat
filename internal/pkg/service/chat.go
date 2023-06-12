package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Chat struct {
	WsUrl string

	tokens map[string]string		//token : login
	activeUsers map[string]string	//token : login

	mtx sync.Mutex
}

func NewChat(wsUrl string) *Chat {
	return &Chat{
		WsUrl:  wsUrl,
		tokens: map[string]string{},
		activeUsers: map[string]string{},}
}

func (c *Chat) GetNewUrl(login string) string {
	token := uuid.New().String()

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.tokens[token] = login

	return c.WsUrl + "?token=" + token
}

func (c *Chat) LoginToken(token string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	login, ok := c.tokens[token]
	if !ok {
		return fmt.Errorf("invalid token %s", token)
	}

	c.activeUsers[token] = login
	delete(c.tokens, token)
	return nil
}

func (c *Chat) LogoutToken(token string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.activeUsers, token)
}

func (c *Chat)GetActiveUsers() []string {
	
	c.mtx.Lock()
	defer c.mtx.Unlock()

	ret := make( []string, 0, len(c.activeUsers))

	for _, v := range c.activeUsers {
		ret = append( ret, v)
	}

	return ret
}
