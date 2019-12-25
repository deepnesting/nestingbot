package delivery

import (
	"testing"

	"github.com/deepnesting/nestingbot/domain/offers"
	"github.com/deepnesting/nestingbot/domain/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestUploadImage(t *testing.T) {
	testTelegramToken := "433942598:AAH9FYcza_C390qmPxNX10Rv24VlS6Jh7lU"
	bot, err := tgbotapi.NewBotAPI(testTelegramToken)
	if err != nil {
		t.Error(err)
	}

	cfg := tgbotapi.NewPhotoUpload(102710272, "/Users/god/Pictures/Screenshots/2.png")
	msg, err := bot.Send(cfg)
	if err != nil {
		t.Error(err)
	}
	if msg.Photo != nil && len(*msg.Photo) > 0 {
		fileID := getMaxFileID(*msg.Photo)
		t.Errorf("%s", fileID)
	}

}

func TestSendFullOffer(t *testing.T) {
	testTelegramToken := "433942598:AAH9FYcza_C390qmPxNX10Rv24VlS6Jh7lU"
	bot, err := tgbotapi.NewBotAPI(testTelegramToken)
	if err != nil {
		t.Error(err)
	}

	testOffer := offers.Offer{
		Text:     "test text",
		Contacts: "адрес: дадада",
		Images:   []string{"AgADAgAD5awxG9KaGUhV-fQbX9x5XhfYug8ABAEAAwIAA3cAAyoqBwABFgQ", "AgADAgAD5qwxG9KaGUifdJSdrQvmqspYzQ8ABAEAAwIAA3gABJUAAhYE"},
	}
	testUser := user.User{
		TID:      102710272,
		Username: "zhuha",
	}

	err = SendFullOfferToUser(bot, testUser.TID, &testOffer, &testUser)
	if err != nil {
		t.Error(err)
	}
}
