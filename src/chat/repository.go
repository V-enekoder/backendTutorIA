package chat

import (
	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"
)

func CreateChatRepository(chat schema.Chat) error {
	db := config.DB

	if err := db.Create(&chat).Error; err != nil {
		return err
	}
	return nil
}

func GetChatsByUserIDAndContextIDRepository(userID uint, contextID uint) ([]schema.Chat, error) {
	db := config.DB
	var chats []schema.Chat

	err := db.Where("user_id = ? AND context_id = ?", userID, contextID).
		Order("created_at asc").
		Find(&chats).Error

	if err != nil {
		return nil, err
	}

	return chats, nil
}
