package service

import (
	"context"
	"errors"
	"time"

	"github.com/fillybodyknow/websocket-chat-app/internal/repository"
	"github.com/fillybodyknow/websocket-chat-app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatServiceInterfact interface {
	SaveMessage(msg *models.ChatMessage) error
}

type ChatService struct {
	ChatRepository repository.ChatRepositoryInterfact
}

func NewChatService(chatRepository *repository.ChatRepositoryInterfact) *ChatService {
	return &ChatService{ChatRepository: *chatRepository}
}

func (s *ChatService) SaveMessage(msg *models.ChatMessage, roomIDStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roomID, err := primitive.ObjectIDFromHex(roomIDStr)
	if err != nil {
		return errors.New("invalid room ID format")
	}

	_, err = s.ChatRepository.FindRoomChatByID(ctx, roomID)
	if err != nil {
		return errors.New("chat room not found")
	}

	return s.ChatRepository.InsertChat(ctx, msg, roomID)
}

func (s *ChatService) CreateRoom(name string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.ChatRepository.InsertRoom(ctx, name)
}
