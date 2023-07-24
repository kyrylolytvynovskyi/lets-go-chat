package restapi2

import (
	//"encoding/json"
	"log"
	"net/http"
	"runtime/trace"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

/*
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	//	EnableCompression: true,
	//	CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	},
}*/

func (srv *server) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("wsEndpoint called")
	token := r.URL.Query().Get("token")
	log.Println("Received token: " + token)

	//check token here
	login, err := srv.chatService.ValidateAndRemoveToken(token)
	if err != nil {
		log.Println("Token validation error:", err)
		return
	}

	wsconn, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade err:", err)
		return
	}
	defer wsconn.Close()

	if err := srv.chatService.LoginToken(token, login, wsconn); err != nil {
		log.Println("Login err:", err)
		return
	}
	defer srv.chatService.LogoutToken(token)
	srv.chanWsConns <- wsconn

	log.Println("processing messages")
	srv.processMessages(login, wsconn)
}

func (srv *server) processMessages(login string, wsconn *websocket.Conn) {
	for {
		messageType, buf, err := wsconn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("processMessages(%v): %v, %v\n", login, messageType, string(buf))

		msg := login + ": " + string(buf)

		srv.chanMessages <- msg

	}
}

func (srv *server) getActiveUsers(w http.ResponseWriter, r *http.Request) {
	defer trace.StartRegion(srv.ctx, "getActiveUsers").End()
	activeUsers := srv.chatService.GetActiveUsers(srv.ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activeUsers)
}
