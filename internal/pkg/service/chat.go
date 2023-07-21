package service

import (
	"context"
	"fmt"
	"log"
	"runtime/trace"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Chat struct {
	WsUrl string

	tokens      map[string]string         //token : login
	activeUsers map[string]connectionInfo //token : login

	mtx sync.Mutex
	wg  sync.WaitGroup
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
	log.Printf("ProcessMessage: %s", msg)

}

func (c *Chat) Run(ctx context.Context, messages <-chan string, wsconns <-chan *websocket.Conn) {

	log.Printf("Start working thread")

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		msgStore := []string{} //msgStore is local variable and do not need any synchronization
		//processing all messages sends in a single thread helps to guarantee proper delivery order

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messages:
				msgStore = append(msgStore, msg)
				c.broadcast(ctx, msg)
			case wsconn := <-wsconns:
				c.sendOldMessages(ctx, msgStore, wsconn)
			}
		}
	}()
}

func (c *Chat) Wait() {
	log.Printf("Wait working thread to gracefully stop")
	c.wg.Wait()
	log.Printf("Working thread gracefully stopped")
}

func (c *Chat) broadcast(ctx context.Context, msg string) {
	defer trace.StartRegion(ctx, "broadcast").End()

	c.mtx.Lock()
	defer c.mtx.Unlock()

	var wg sync.WaitGroup
	defer wg.Wait() //we need to wait until broadcast is finished to maintain proper order of messages

	for _, v := range c.activeUsers {
		//send message to each user in separate async goroutine, then wait until all of them finishes
		wg.Add(1)
		go func(wsconn *websocket.Conn, login string) {
			defer wg.Done()
			log.Printf("broadcasting message %s to conn %s", msg, login)

			if err := wsconn.WriteMessage(1, []byte(msg)); err != nil {
				log.Println(err)
			}
		}(v.wsconn, v.login)
	}

}

func (c *Chat) sendOldMessages(ctx context.Context, msgStore []string, wsconn *websocket.Conn) {
	defer trace.StartRegion(ctx, "sendOldMessages").End()

	for _, msg := range msgStore {
		log.Printf("sending old message %s", msg)

		if err := wsconn.WriteMessage(1, []byte(msg)); err != nil {
			log.Println(err)
		}
	}
}
