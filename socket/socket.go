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
	Token        string
	SendNewOrder func(o *gorm.Order)
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
			if gorm.ValidateAdminSocketToken(token) != nil {
				logger.Info("token validation failed")
				return "error"
			}

			// register the client as admin
			c := client{
				Token: token,
			}
			adminClients = append(adminClients, &c)
			logger.Info("socket id", so.Id(), "registered as admin.")

			// enable admin actions
			so.On("qfeed-commit-queue", handleCommitQueue)
			so.On("qfeed-commit-prep", handleCommitPrep)
			so.On("qfeed-commit-ship", handleCommitShip)

			c.SendNewOrder = func(o *gorm.Order) {
				so.Emit("qfeed-new-order", o)
			}

			return "okbro"
		})
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

// NewOrderPlaced is used after user placed an order. Will send notification to the admin queue
// and the order tracker
func NewOrderPlaced(o *gorm.Order) {
	// send order to the admin queue
	for _, c := range adminClients {
		c.SendNewOrder(o)
	}
	// TODO: send something to the chatbot
	return
}
