package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/fillybodyknow/websocket-chat-app/internal/service"
	"github.com/fillybodyknow/websocket-chat-app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	ChatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{ChatService: chatService}
}

func (h *ChatHandler) WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	roomID := c.Query("room_id")
	sender := c.Query("sender")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Failed to read message:", err)
			break
		}

		log.Printf("📩 Received: %s\n", message)

		chat := &models.ChatMessage{Sender: sender, Message: string(message), CreatedAt: time.Now()}

		if err := h.ChatService.SaveMessage(chat, roomID); err != nil {
			log.Println("❌ Failed to save message:", err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte("✅ Message received")); err != nil {
			log.Println("❌ Failed to write message:", err)
			break
		}
	}
}

func (h *ChatHandler) CreateRoomHandler(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID, err := h.ChatService.CreateRoom(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างห้องแชทไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "สร้างห้องแชทสำเร็จ", "room_id": roomID.Hex()})
}
