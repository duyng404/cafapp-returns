package socket

import (
	"cafapp-returns/logger"
	"github.com/googollee/go-socket.io"
)

//Server export
var Server, err = socketio.NewServer(nil)

func init() {
	if err != nil {
		logger.Fatal(err)
	}
	Server.On("connection", func(so socketio.Socket) {
		logger.Info("on connection")
		so.Join("chat")
		so.On("createMessage", func(message map[string]string) {
			//create a new message map
			newMessage := make(map[string]string)
			newMessage["user"] = message["user"]
			newMessage["content"] = message["content"]
			Server.BroadcastTo("chat", "newMessage", newMessage)
		})
		so.On("disconnection", func() {
			logger.Info("on disconnect")
		})
	})
	Server.On("error", func(so socketio.Socket, err error) {
		logger.Info("error:", err)
	})
}
