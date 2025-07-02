package router

import (
	"github.com/fillybodyknow/websocket-chat-app/internal/handler"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	ChatHandler *handler.ChatHandler
}

func NewChatRouter(chatHandler *handler.ChatHandler) *ChatRouter {
	return &ChatRouter{ChatHandler: chatHandler}
}

func (r *ChatRouter) ChatRoutes(router *gin.RouterGroup) {
	router.GET("/ws", r.ChatHandler.WebsocketHandler)
	router.POST("/create_room", r.ChatHandler.CreateRoomHandler)
}
