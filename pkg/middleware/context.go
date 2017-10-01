package middleware

import (
	"github.com/deepnesting/nestingbot/pkg/models"

	"github.com/zhuharev/tamework"
)

func MakeUserFromUpdate(update *tamework.Update) *models.User {
	user := models.User{
		TelegramID: update.ChatID(),
		FirstName:  update.FirstName(),
		LastName:   update.LastName(),

		Username: update.Username(),
	}
	return &user
}
