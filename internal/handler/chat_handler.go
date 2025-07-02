package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/fillybodyknow/websocket-chat-app/internal/hub"
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
	Hub         *hub.Hub
}

func NewChatHandler(chatService *service.ChatService, hubInstance *hub.Hub) *ChatHandler {
	return &ChatHandler{
		ChatService: chatService,
		Hub:         hubInstance,
	}
}

func handleRead(client *hub.Client, sender string, service *service.ChatService, hubInstance *hub.Hub) {
	defer func() {
		hubInstance.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("❌ Read error:", err)
			break
		}

		service.SaveMessage(&models.ChatMessage{
			Sender:    sender,
			Message:   string(msg),
			CreatedAt: time.Now(),
		}, client.RoomID)

		broadcastMsg := []byte(sender + "|" + string(msg))

		hubInstance.Broadcast <- hub.Message{
			RoomID:    client.RoomID,
			Data:      broadcastMsg,
			SenderRef: client,
		}

	}
}

func handleWrite(client *hub.Client) {
	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("❌ Write error:", err)
			break
		}
	}
}

func (h *ChatHandler) WebsocketHandler(c *gin.Context) {
	roomID := c.Query("room_id")
	sender := c.Query("sender")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade connection:", err)
		return
	}

	client := &hub.Client{
		Conn:   conn,
		Send:   make(chan []byte),
		RoomID: roomID,
	}

	h.Hub.Register <- client

	// Read & write goroutines
	go handleRead(client, sender, h.ChatService, h.Hub)
	go handleWrite(client)
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

func (h *ChatHandler) GetHistoryChatHandler(c *gin.Context) {
	roomID := c.Query("room_id")
	chatMessages, err := h.ChatService.GetHistoryChat(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการดึงข้อมูล"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat_messages": chatMessages})
}
