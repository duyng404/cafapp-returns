package socket

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"

	"github.com/googollee/go-socket.io"
)

// private vars
var adminClients []*client
var userClients []*client
var server *socketio.Server

// client
type client struct {
	ID                 string
	User               *gorm.User
	Token              string
	SendNewOrder       func(o *gorm.Order)
	UpdateCurrentQueue func(orders []*gorm.Order)
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

		// handle registration for admin
		so.On("register", func(token string) string {
			// admin registration: add the socket to the adminClients list, and return ack
			user, err := gorm.ValidateAdminSocketToken(token)
			if err != nil {
				logger.Info("token validation failed")
				return "error"
			}

			// register the client as admin
			c := client{
				Token: token,
				ID:    so.Id(),
				User:  user,
			}
			adminClients = append(adminClients, &c)
			logger.Info("socket id", so.Id(), "registered as admin.")

			// enable admin actions
			so.On("qfeed-commit-queue", func(committed []int) string {
				handleCommit("qfeed-commit-queue", committed)
				return "okbro"
			})
			so.On("qfeed-commit-prep", func(committed []int) string {
				handleCommit("qfeed-commit-prep", committed)
				return "okbro"
			})
			so.On("qfeed-commit-ship", func(committed []int) string {
				handleCommit("qfeed-commit-ship", committed)
				return "okbro"
			})
			so.On("qfeed-commit-approach", func(committed []int) string {
				handleCommit("qfeed-commit-approach", committed)
				return "okbro"
			})

			c.SendNewOrder = func(o *gorm.Order) {
				so.Emit("qfeed-new-order", o)
			}

			c.UpdateCurrentQueue = func(orders []*gorm.Order) {
				so.Emit("qfeed-update-queue", orders)
			}

			return "okbro"
		})

		// handle registration for users
		so.On("register-user", func(token string) string {
			// user registration: add the socket to the userClients list, and return ack
			user, err := gorm.ValidateUserSocketToken(token)
			if err != nil {
				logger.Info("token validation failed")
				return "error"
			}

			// register the client as user
			c := client{
				Token: token,
				ID:    so.Id(),
				User:  user,
			}
			userClients = append(userClients, &c)
			logger.Info("socket id", so.Id(), "registered as user.")

			// enable user action
			so.On("chatbot-request", func(request string) string {
				response := c.handleChatbotRequest(request)
				so.Emit("chatbot-response", response)
				return "okbro"
			})

			return "okbro"
		})

		so.On("disconnection", func() {
			// deregister, remove them from admin list
			for i, v := range adminClients {
				if v.ID == so.Id() {
					adminClients = append(adminClients[:i], adminClients[i+1:]...)
					logger.Info("admin socket id", so.Id(), "disconnected.")
					break
				}
			}
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
		logger.Info("sending new order to user", c.User.GusUsername)
		c.SendNewOrder(o)
	}
	// TODO: send something to the chatbot
	return
}

// when one admin changes something, update it for every other admins currently connected.
func updateQueueForEveryone(orders []*gorm.Order) {
	for _, c := range adminClients {
		logger.Info("sending update queue to user", c.User.GusUsername)
		c.UpdateCurrentQueue(orders)
	}
}
