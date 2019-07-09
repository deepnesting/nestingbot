package delivery

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deepnesting/nestingbot/domain/offers"
	"github.com/deepnesting/nestingbot/domain/user"
	"github.com/deepnesting/nestingbot/pkg/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhuharev/tamework"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
)

var (
	_id uint64 = 1000
)

func kbTags(offer *offers.Offer, ctx *tamework.Context) {
	ctx.Keyboard.Reset()
	var btns = []string{
		buttons.TagSearchNester,
		buttons.TagSearchCompanion,
		buttons.TagReplace,
		buttons.TagHookUp,
		buttons.TagSearchNest,
	}
	for _, btn := range btns {
		var found bool
		for _, t := range offer.Tags {
			if t == btn {
				found = true
			}
		}
		if found {
			ctx.Keyboard.AddCallbackButton("✅ "+btn, "toggletag:"+fmt.Sprint(offer.ID)+":"+btn)
			ctx.Keyboard.AddCallbackButton("", "")
		} else {
			ctx.Keyboard.AddCallbackButton("✖️ "+btn, "toggletag:"+fmt.Sprint(offer.ID)+":"+btn)
			ctx.Keyboard.AddCallbackButton("", "")
		}
	}
}

func getMaxFileID(ps []tgbotapi.PhotoSize) string {
	var max = 0
	var maxID string
	for _, s := range ps {
		if s.Width > max {
			max = s.Width
			maxID = s.FileID
		}
	}
	return maxID
}

func getFileID(update tamework.Update) (string, bool) {
	if update.Message == nil {
		return "", false
	}
	if update.Message.Photo != nil && len(*update.Message.Photo) > 0 {
		if id := getMaxFileID(*update.Message.Photo); id != "" {
			return id, true
		}
	}
	if update.Message.Document != nil && update.Message.Document.Thumbnail != nil {
		return update.Message.Document.FileID, true
	}
	return "", false
}

// MakeCreate returns tamework handler
func MakeCreate(typ string, admins []int64, offerRepo offers.Repository) tamework.HandleFunc {
	return func(ctx *tamework.Context) {
		var (
			text     string
			images   []string
			contacts string
		)
		ctx.Keyboard.AddReplyButton(buttons.CancelButton)
		ctx.Send("Шаг 1/4. Пришлите текст объявления.")
		for {
			update, done := ctx.Wait(buttons.CancelButton, time.Minute*5)
			if !done {
				setDefaultKeyboard(ctx)
				ctx.Send("Создание объявления прекращено, попробуйте ещё раз")
				return
			}
			text = update.Text()
			if len(text) < 10 {
				ctx.Send("Текст объявления слишком короткий")
				continue
			} else {
				break
			}
		}
		ctx.Keyboard.AddReplyButton(buttons.CancelButton)
		ctx.Keyboard.AddReplyButton(buttons.NextButton)
		ctx.Send("Шаг 2/4. Пришлите фотографии. Дождитесь пока фотографии загрузятся и нажмите кнопку <Далее>.")
		for {
			update, done := ctx.Wait(buttons.CancelButton, time.Minute*5)
			if !done {
				setDefaultKeyboard(ctx)
				ctx.Send("Создание объявления прекращено, попробуйте ещё раз")
				return
			}
			if update.Text() == buttons.NextButton {
				if len(images) == 0 {
					ctx.Send("Нужна хотя бы одна картинка")
					continue
				}
				break
			}
			if id, ok := getFileID(update); ok {
				images = append(images, id)
				ctx.Send("Фото загружено, загрузите ещё или нажмите Далее")
			} else {
				ctx.Send("Нужно послать фотку")
			}
		}

		ctx.Keyboard.AddReplyButton(buttons.CancelButton)
		ctx.Send("Шаг 3/4. Пришлите/напишите свои контактные данные.")
		for {
			update, done := ctx.Wait(buttons.CancelButton, time.Minute*5)
			if !done {
				setDefaultKeyboard(ctx)
				ctx.Send("Создание объявления прекращено, попробуйте ещё раз")
				return
			}
			if update.Text() == "" {
				ctx.Send("Контакты не могут быть пустыми")
				continue
			}
			contacts = update.Text()
			break
		}

		id, err := offerRepo.Create(int(ctx.UserID), typ, text, images, contacts)
		if err != nil {
			log.Error("adsa", rz.Err(err))
		}

		ctx.Keyboard.AddReplyButton(buttons.CancelButton)
		ctx.Keyboard.AddReplyButton(buttons.NextButton)

		ctx.Send("Последний шаг) Добавьте тэги и нажмите 'Далее'")

		offer, err := offerRepo.GetByID(id)
		if err != nil {
			log.Error("adsa", rz.Err(err))
		}

		userFace := ctx.Data["user"]
		user, ok := userFace.(*user.User)
		if !ok {
			//???

		}

		ctx.Keyboard = tamework.NewKeyboard(nil)

		renderedOffer, _ := offers.FormatMarkdown(*offer, user.Username)
		kbTags(offer, ctx)
		msg, err := ctx.Markdown(renderedOffer)
		if err != nil {
			log.Error("err send msg", rz.Err(err), rz.String("body", msg.Text))
		}

		for {
			update, done := ctx.Wait(buttons.CancelButton)
			if !done {
				ctx.Send("создание объявление отменено")
				return
			}

			txt := update.Text()
			log.Debug("update here", rz.String("text", txt))
			if strings.HasPrefix(txt, "toggletag:") {
				arr := strings.SplitN(txt, ":", 3)
				if len(arr) != 3 {
					log.Error("split", rz.String("text", txt))
					return
				}
				offerID, err := strconv.Atoi(arr[1])
				if err != nil {
					log.Error("parse offer id", rz.Err(err))
					return
				}
				if offerID != id {
					log.Error("bad offer id", rz.Int("expected", id),
						rz.Int("got", offerID))
				}
				log.Debug("toggle tag", rz.String("tag", arr[2]), rz.Int("offer_id", id))
				err = offerRepo.ToggleTag(id, arr[2])
				if err != nil {
					log.Error("toggle tag in db", rz.Err(err))
					return
				}
				oldTags := offer.Tags
				offer, err = offerRepo.GetByID(offer.ID)
				if err != nil {
					log.Error("get offer from db", rz.Err(err))
					return
				}
				log.Debug("toggled tags",
					rz.Strings("old_tags", oldTags),
					rz.Strings("new_tags", offer.Tags),
				)
				renderedOffer, _ := offers.FormatMarkdown(*offer, user.Username)
				log.Debug("rendered", rz.String("text", renderedOffer))
				// update msg text
				cnf := tgbotapi.NewEditMessageText(ctx.ChatID, msg.MessageID, renderedOffer)
				cnf.ParseMode = tgbotapi.ModeMarkdown
				kbTags(offer, ctx)
				if kb, ok := ctx.Keyboard.Markup().(tgbotapi.InlineKeyboardMarkup); ok {
					cnf.ReplyMarkup = &kb
				}
				_, err = ctx.BotAPI().Send(cnf)
				if err != nil {
					log.Error("err send msg", rz.Err(err), rz.String("body", msg.Text))
				}
				continue
			} else {
				log.Debug("break", rz.String("text", update.Text()))

				setDefaultKeyboard(ctx)
				ctx.Send("Объявление отправлено на модерацию")

				break
			}
		}
		ctx.Keyboard.Reset()
		ctx.Keyboard.AddURLButton("Оплатить", getPaymentURL(uint64(offer.ID)))
		ctx.Markdown("Объявление будет опубликовано после оплаты")

		// renderedOffer, _ = offers.FormatMarkdown(*offer, user.Username)
		// for _, admin := range admins {
		// 	ctx.Keyboard.AddCallbackButton("Меню", fmt.Sprintf("showmenu:%d", offer.ID))
		// 	ctx.MarkdownTo(admin, renderedOffer)
		// }
	}
}

func getPaymentURL(id uint64) string {
	uri := "https://money.yandex.ru/quickpay/shop-widget?writer=seller&targets-hint=&button-text=11&hint=&quickpay=shop&payment-type-choice=on"
	uri += "&account=" + os.Getenv("YANDEX_MONEY_ACCOUNT")
	uri += "&targets=Объявление в Телеграмм-канале"
	uri += "&default-sum=39"
	uri += fmt.Sprintf("&label=u:%d", id)
	uri += "&successURL=https://t.me/ugnestbot"
	return uri
}

func setDefaultKeyboard(c *tamework.Context) {
	c.NewKeyboard(buttons.Menu)
	c.Keyboard.SetRowLen(2)
	c.Keyboard.SetType(tamework.KeyboardReply)
}
