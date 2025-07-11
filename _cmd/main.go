package main

import (
	"log"
	"os"

	"github.com/fillybodyknow/websocket-chat-app/config"
	"github.com/fillybodyknow/websocket-chat-app/internal/handler"
	"github.com/fillybodyknow/websocket-chat-app/internal/hub"
	"github.com/fillybodyknow/websocket-chat-app/internal/repository"
	"github.com/fillybodyknow/websocket-chat-app/internal/router"
	"github.com/fillybodyknow/websocket-chat-app/internal/service"
	"github.com/gin-contrib/cors"
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
	port := os.Getenv("PORT")
	NGROK8080 := os.Getenv("NGROK_8080")
	NGROK5500 := os.Getenv("NGROK_5500")

	hubInstance := hub.NewHub()
	go hubInstance.Run()

	ChatCollection := db.Database(DBname).Collection("chat_rooms")
	ChatRepo := repository.NewChatRepository(ChatCollection)
	ChatService := service.NewChatService(ChatRepo)
	ChatHandler := handler.NewChatHandler(ChatService, hubInstance)
	ChatRouter := router.NewChatRouter(ChatHandler)
	gin.SetMode(gin.DebugMode)
	Router := gin.Default()
	Router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			NGROK8080,
			NGROK5500,
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "ngrok-skip-browser-warning"},
		AllowCredentials: true,
	}))

	API := Router.Group("/api")
	Chat := API.Group("/chat_rooms")
	ChatRouter.ChatRoutes(Chat)

	Router.Run(":" + port)
}
