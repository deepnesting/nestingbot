package subscriptions

import (
	"fmt"

	"github.com/deepnesting/nestingbot/pkg/buttons"
	"github.com/deepnesting/nestingbot/pkg/models"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/zhuharev/tamework"
)

func Subscriptions(c *tamework.Context) {
	c.NewKeyboard(buttons.FiltersMenu)
	c.Send("Какие объявления вы хотите получать?")
}

func SubscriptionsOwner(c *tamework.Context) {
	c.NewKeyboard(buttons.FiltersOwnerMenu)
	c.Send("Какие объявления вы хотите получать?")
}

func MakeTogleSubscribe(cat string) tamework.HandleFunc {
	return func(c *tamework.Context) {
		color.Cyan("Subscribe %s %d", cat, c.ChatID)
		_, err := models.UserGetOrCreate(c.ChatID)
		if err != nil {
			color.Red("%s", err)
			return
		}
		channels, err := models.GetSubscriptions(c.ChatID)
		if err != nil {
			color.Red("%s", err)
			return
		}
		var alreadySubscribed bool
		for _, v := range channels {
			if v == cat {
				alreadySubscribed = true
			}
		}

		if alreadySubscribed {
			c.Keyboard.AddCallbackButton("Отписаться", "un"+cat)
			c.Markdown(fmt.Sprintf("Вы уже подписаны на категорию *%s*", cat))
			return
		} else {
			err = models.Subscribe(c.ChatID, cat)
			if err != nil {
				color.Red("%s", err)
				return
			}
		}

		cnt, err := models.GetSubscribersCount(cat)
		if err != nil {
			color.Red("%s", err)
			return
		}

		c.Keyboard.AddCallbackButton("Отписаться", "un"+cat)
		c.Send(fmt.Sprintf("Вы успешно подписались (всего подписчиков: %d)", cnt))
	}
}

func MakeTogleUnSubscribe(cat string) tamework.HandleFunc {
	return func(c *tamework.Context) {
		color.Cyan("Unsubscribe %s %d", cat, c.ChatID)
		_, err := models.UserGetOrCreate(c.ChatID)
		if err != nil {
			color.Red("%s", err)
			return
		}

		err = models.Unsibscribe(c.ChatID, cat)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.Markdown(fmt.Sprintf("Вы ещё не подписаны на категорию *%s*", cat))
				return
			}
			color.Red("%s", err)
			return
		}
		c.Send("Вы успешно отписались")
	}
}
