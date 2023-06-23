package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Chat struct {
	WsUrl string

	tokens      map[string]string         //token : login
	activeUsers map[string]connectionInfo //token : login

	mtx sync.Mutex
}

type connectionInfo struct {
	login  string
	wsconn *websocket.Conn
}

func NewChat(wsUrl string) *Chat {
	return &Chat{
		WsUrl:       wsUrl,
		tokens:      map[string]string{},
		activeUsers: map[string]connectionInfo{}}
}

func (c *Chat) GetNewUrl(login string) string {
	token := uuid.New().String()

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.tokens[token] = login

	return c.WsUrl + "?token=" + token
}

// returns login for given token
func (c *Chat) ValidateAndRemoveToken(token string) (string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	login, ok := c.tokens[token]
	if !ok {
		return "", fmt.Errorf("invalid token %s", token)
	}

	delete(c.tokens, token)

	return login, nil
}

func (c *Chat) LoginToken(token string, login string, wsconn *websocket.Conn) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.activeUsers[token] = connectionInfo{login, wsconn}
	return nil
}

func (c *Chat) LogoutToken(token string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	_, ok := c.activeUsers[token]
	if !ok {
		return fmt.Errorf("invalid token %s", token)
	}

	delete(c.activeUsers, token)
	return nil
}

func (c *Chat) GetActiveUsers() []string {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	ret := make([]string, 0, len(c.activeUsers))

	for _, v := range c.activeUsers {
		ret = append(ret, v.login)
	}

	return ret
}

func (c *Chat) ProcessMessage(messageType int, msg string) {

}

func (c *Chat) Run(ctx context.Context, messages <-chan string, wsconns <-chan *websocket.Conn) {

	go func() {
		msgStore := []string{}

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messages:
				msgStore = append(msgStore, msg)
				c.broadcast(msg)
			case wsconn := <-wsconns:
				c.sendOldMessages(msgStore, wsconn)
			}
		}
	}()
}

func (c *Chat) broadcast(msg string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, v := range c.activeUsers {

		log.Printf("broadcasting message %s to conn %s", msg, v.login)
	}
}

func (c *Chat) sendOldMessages(msgStore []string, wsconn *websocket.Conn) {

	for _, msg := range msgStore {
		log.Printf("sending message %s", msg)
	}
}
