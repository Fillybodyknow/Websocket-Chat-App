package repository

import (
	"context"
	"errors"

	"github.com/fillybodyknow/websocket-chat-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepositoryInterfact interface {
	InsertChat(ctx context.Context, chat *models.ChatMessage, roomID primitive.ObjectID) error
	FindRoomChatByID(ctx context.Context, roomID primitive.ObjectID) (models.ChatRoom, error)
	InsertRoom(ctx context.Context, name string) (primitive.ObjectID, error)
}

type ChatRepository struct {
	ChatCollection *mongo.Collection
}

func NewChatRepository(chatCollection *mongo.Collection) *ChatRepository {
	return &ChatRepository{
		ChatCollection: chatCollection,
	}
}

func (r *ChatRepository) InsertChat(ctx context.Context, chat *models.ChatMessage, roomID primitive.ObjectID) error {
	_, err := r.ChatCollection.UpdateOne(
		ctx,
		bson.M{"_id": roomID},
		bson.M{"$push": bson.M{"chat_messages": chat}},
	)
	return err
}

func (r *ChatRepository) FindRoomChatByID(ctx context.Context, roomID primitive.ObjectID) (models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.ChatCollection.FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)
	return room, err
}

func (r *ChatRepository) InsertRoom(ctx context.Context, name string) (primitive.ObjectID, error) {
	newRoom := models.ChatRoom{
		ID:           primitive.NewObjectID(),
		Name:         name,
		ChatMessages: []models.ChatMessage{},
	}

	res, err := r.ChatCollection.InsertOne(ctx, newRoom)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to cast inserted ID")
	}

	return insertedID, nil
}
