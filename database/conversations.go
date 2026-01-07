package database

import (
	"chatty/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetConversation(user1, user2 string, database *mongo.Database) *mongo.Collection {
	conversationHash := utils.CommutativeUUIDHashFromString(user1, user2)
	conversation := database.Collection(conversationHash)
	return conversation
}
