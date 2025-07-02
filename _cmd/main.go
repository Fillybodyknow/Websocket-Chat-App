package main

import (
	"log"
	"os"

	"github.com/fillybodyknow/websocket-chat-app/config"
	"github.com/fillybodyknow/websocket-chat-app/internal/handler"
	"github.com/fillybodyknow/websocket-chat-app/internal/repository"
	"github.com/fillybodyknow/websocket-chat-app/internal/router"
	"github.com/fillybodyknow/websocket-chat-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := config.ConnectMongoDB()
	DBname := os.Getenv("DB_NAME")

	ChatCollection := db.Database(DBname).Collection("chat_rooms")
	ChatRepo := repository.NewChatRepository(ChatCollection)
	ChatService := service.NewChatService(&ChatRepo)
	ChatHandler := handler.NewChatHandler(ChatService)
	ChatRouter := router.NewChatRouter(ChatHandler)
	Router := gin.Default()
	Chat := Router.Group("/chat")
	ChatRouter.ChatRoutes(Chat)

	port := os.Getenv("PORT")
	Router.Run(":" + port)
}
