package socket

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"

	"github.com/googollee/go-socket.io"
)

// private vars
var adminClients []*client
var server *socketio.Server

// client
type client struct {
	Token  string
	Socket *socketio.Socket
}

// socketMessage
type socketMessage struct {
	Msg string `json:"msg"`
}

func init() {
	// init the server
	logger.Info("initializing socketio server")
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		logger.Fatal(err)
	}

	// lifecycle of a connection, aka, socket's routes
	server.On("connection", func(so socketio.Socket) {
		logger.Info("socket connection from", so.Id())

		// handle registration
		so.On("register", func(token string) string {
			// admin registration: add the socket to the adminClients list, and return ack
			if gorm.ValidateAdminSocketToken(token) == nil {
				c := client{
					Token:  token,
					Socket: &so,
				}
				adminClients = append(adminClients, &c)
				logger.Info("socket id", so.Id(), "registered as admin.")
				return "okbro"
			}
			logger.Info("token validation failed")
			return "error"
		})

		// handle all other stuff
		so.On("qfeed-commit-queue", commitQueue)

		so.On("disconnection", func() {
			logger.Info("socket disconnected", so.Id())
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		logger.Info("error:", err)
	})
}

// GetServer get the current server
func GetServer() *socketio.Server {
	return server
}

func commitQueue() string {
	return ""
}
