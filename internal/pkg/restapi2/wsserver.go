package restapi2

import (
	"encoding/json"
	"log"
	"net/http"

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
	if err := srv.chatService.LoginToken(token); err != nil {
		log.Println("Token validation error:", err)
		return
	}
	defer srv.chatService.LogoutToken(token)

	wsconn, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade err:", err)
		return
	}
	defer wsconn.Close()

	log.Println("processing messages")
	srv.processMessages(wsconn)

}

func (srv *server) processMessages(wsconn *websocket.Conn) {
	for {
		messageType, buf, err := wsconn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(buf))

		if err := wsconn.WriteMessage(messageType, buf); err != nil {
			log.Println(err)
			return
		}
	}
}

func (srv *server) getActiveUsers(w http.ResponseWriter, r *http.Request) {

	activeUsers := srv.chatService.GetActiveUsers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activeUsers)
}
