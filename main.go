package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"

	smmmodels "pure/smm/boards/models"

	"github.com/Unknwon/com"
	"github.com/deepnesting/nestingbot/pkg/binlog"
	"github.com/deepnesting/nestingbot/pkg/buttons"
	"github.com/deepnesting/nestingbot/pkg/models"
	"github.com/deepnesting/nestingbot/pkg/setting"
	"github.com/deepnesting/nestingbot/routers/subscriptions"
	"github.com/fatih/color"
	"github.com/go-macaron/binding"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	dry "github.com/ungerik/go-dry"
	"github.com/zhuharev/tamework"
)

var (
	tw *tamework.Tamework
)

func main() {
	err := setting.NewContext()
	if err != nil {
		log.Fatalf("Failed setting initialization, %s", err)
	}

	err = models.NewContext()
	if err != nil {
		log.Fatalf("Failed boltdb initialization, %s", err)
	}
	err = binlog.NewContext()
	if err != nil {
		log.Fatalf("Failed binlog initialization, %s", err)
	}

	tw, err = tamework.New(setting.App.Telegram.Token)
	if err != nil {
		log.Fatalf("Failed telegram bot initialization, %s", err)
	}

	tw.Text("/menu", Greeting)
	//alias
	tw.RegistreMethod("/menu", "/start")
	tw.RegistreMethod("/menu", buttons.Back)
	tw.RegistreMethod("/menu", buttons.CancelButton)
	//tw.Text("/start", Greeting)
	tw.Text(buttons.SubscriptionsButton, subscriptions.Subscriptions)
	tw.Text(buttons.SubscriptionsOwnerButton, subscriptions.SubscriptionsOwner)

	var subscribeButtons = []string{
		buttons.AboutRent,
		buttons.AboutNeightborg,
		buttons.AboutRentRoom,
	}

	for _, but := range subscribeButtons {
		tw.Text(but, subscriptions.MakeTogleSubscribe(but))
		tw.CallbackQuery("un"+but, subscriptions.MakeTogleUnSubscribe(but))
	}

	var subscribeFinderButtons = []string{
		buttons.AboutFinderRoom,
		buttons.AboutFinder,
	}

	for _, but := range subscribeFinderButtons {
		tw.Text(but, subscriptions.MakeTogleSubscribe(but))
		tw.CallbackQuery("un"+but, subscriptions.MakeTogleUnSubscribe(but))
	}

	tw.Text("/support", Support)
	tw.CallbackQuery("/support", Support)
	tw.RegistreMethod("/support", buttons.HelpButton)

	tw.Text(buttons.AddProposalSearch, Add)
	//tw.RegistreMethod(buttons.AddProposal, buttons.AddProposalSearch)

	tw.Text("/terms", Terms)

	// chat
	tw.Prefix("sup", Sup)

	tw.Prefix("setcat_", SetCat)

	tw.Prefix("upvote_", UpVote)
	tw.Prefix("downvote_", DownVote)
	tw.Prefix("publish_", Publish)

	go tw.Run()

	m := macaron.New()
	m.Use(macaron.Renderer())

	m.Post(fmt.Sprintf("/%s/event", setting.App.Secret), binding.Bind(Message{}), eventHandler)

	m.Run(2018)

}

func Publish(c *tamework.Context) {
	c.Answer("опубликовано")
	msgID := com.StrTo(c.Text).MustInt64()

	var res smmmodels.Messages

	err := dry.FileUnmarshallJSON(fmt.Sprintf("https://smmpolice.ru/api/v1/messages/%d", msgID), &res)
	if err != nil {
		c.Send(err.Error())
		return
	}

	uploadRemoteMessage("@ughome", msgID, res.UserLink(), getVoteKeyboard(msgID, false))
}

func UpVote(c *tamework.Context) {
	vote(c, 1)
}

func DownVote(c *tamework.Context) {
	vote(c, -1)
}

func vote(c *tamework.Context, vote int) {
	user, err := models.UserGetOrCreate(c.UserID)
	if err != nil {
		color.Red("%s", err)
		return
	}
	objID := com.StrTo(c.Text).MustInt64()
	err = models.VotesVote(user.ID, objID, vote)
	if err != nil {
		color.Red("%s", err)
		return
	}
	c.Answer("Голос учтён")
	cnt, err := models.VotesCount(objID, 1)
	if err != nil {
		color.Red("%s", err)
		return
	}
	cntDown, err := models.VotesCount(objID, -1)
	if err != nil {
		color.Red("%s", err)
		return
	}
	isAdmin := c.ChatID == 102710272
	c.EditReplyMurkup(getVoteKeyboard(objID, isAdmin, cnt, cntDown))
}

func SetCat(c *tamework.Context) {
	arr := strings.Split(c.Text, "_")
	if len(arr) != 2 {
		return
	}
	msgID := com.StrTo(arr[0]).MustInt64()
	catID := com.StrTo(arr[1]).MustInt()
	if cap, ok := catWait[msgID]; ok {
		err := broadcast(msgID, catID, cap)
		if err != nil {
			log.Println(err)
		}
		delete(catWait, msgID)
	}
	c.Answer("Сообщение опубликовано")
	c.EditText(c.Update().CallbackQuery.Message.Text)
}

func Sup(c *tamework.Context) {
	userID := com.StrTo(c.Text).MustInt64()
	c.Keyboard.AddReplyButton(buttons.CancelButton)
	c.Send("Введите ответное сообщение:")
	u, done := c.Wait(buttons.CancelButton, time.Second*180)
	if done {
		c.Keyboard.Reset().AddCallbackButton("Ответить", "/support")
		c.SendTo(userID, u.Text())
		setDefaultKeyboard(c)
		c.Send("Сообщение отправлено.")
	}
}

var catWait = map[int64]Message{}

func eventHandler(c *macaron.Context, event Message) {
	if event.IsTest {
		c.JSON(200, "ok")
		return
	}

	if event.Category == 0 {
		kb := tamework.NewKeyboard(nil).AddCallbackButton("О сдаче квартир", fmt.Sprintf("setcat_%d_1", event.ID)).
			AddCallbackButton("").
			AddCallbackButton("О сдаче комнат", fmt.Sprintf("setcat_%d_2", event.ID)).
			AddCallbackButton("").
			AddCallbackButton("О поиске соседа", fmt.Sprintf("setcat_%d_3", event.ID)).
			AddCallbackButton("").
			AddCallbackButton("О поиске комнат", fmt.Sprintf("setcat_%d_4", event.ID)).
			AddCallbackButton("").
			AddCallbackButton("О поиск квартир", fmt.Sprintf("setcat_%d_5", event.ID))
		if event.HasPhoto {
			uploadRemoteMessage(setting.App.Telegram.Admin, event.ID, event.Contact, kb)

		} else {
			msg := tgbotapi.NewMessage(setting.App.Telegram.Admin,
				fmt.Sprintf("%s\n\n%s", event.Body, event.Contact))
			msg.ReplyMarkup = kb.Markup()
			_, err := tw.Bot().Send(msg)
			if err != nil {
				color.Red("%s", err)
			}
		}
		catWait[event.ID] = event
		return
	}
	err := broadcast(event.ID, event.Category, event)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, "ok")
}

func getVoteKeyboard(msgID int64, isAdmin bool, votes ...int) *tamework.Keyboard {
	kb := tamework.NewKeyboard(nil)
	up := "👍"
	if len(votes) > 0 && votes[0] != 0 {
		up += fmt.Sprintf(" %d", votes[0])
	}
	down := "👎"
	if len(votes) > 1 && votes[1] != 0 {
		down += fmt.Sprintf(" %d", votes[1])
	}
	kb.AddCallbackButton(up, "upvote_"+fmt.Sprint(msgID)).
		AddCallbackButton(down, "downvote_"+fmt.Sprint(msgID))
	if isAdmin {
		kb.AddCallbackButton("в канал "+fmt.Sprint(msgID), "publish_"+fmt.Sprint(msgID))
	}
	return kb
}

func broadcast(msgID int64, catID int, event Message) (err error) {
	var (
		subs      []int64
		channelID = ""
	)

	switch catID {
	//AboutRent
	case 1:
		channelID = buttons.AboutRent
	case 2:
		channelID = buttons.AboutRentRoom
	case 3:
		channelID = buttons.AboutNeightborg
	case 4:
		channelID = buttons.AboutFinderRoom
	case 5:
		channelID = buttons.AboutFinder
	default:
		return
	}

	subs, err = models.GetSubscribers(channelID)
	if err != nil {
		return
	}

	color.Green("%v", subs)
	var (
		photoID string
	)

	for _, subID := range subs {
		kb := getVoteKeyboard(msgID, false) //

		if photoID == "" && event.HasPhoto {
			photoID, err = uploadRemoteMessage(subID, msgID, event.Contact, kb)
			if err != nil {
				color.Red("%s", err)
				continue
			}
			continue
		}
		if event.HasPhoto {
			msg := tgbotapi.NewPhotoShare(subID, photoID)
			msg.Caption = event.Contact
			msg.ReplyMarkup = kb.Markup()
			_, err = tw.Bot().Send(msg)
			if err != nil {
				color.Red("%s", err)
				continue
			}
		} else {
			msg := tgbotapi.NewMessage(subID,
				fmt.Sprintf("%s\n\n%s", event.Body, event.Contact))
			msg.ReplyMarkup = kb.Markup()
			_, err = tw.Bot().Send(msg)
			if err != nil {
				color.Red("%s", err)
				continue
			}
		}

	}
	return
}

func uploadRemoteMessage(userID interface{}, msgID int64, caption string, kbs ...*tamework.Keyboard) (photoID string, err error) {
	var imageURI = "https://smmpolice.ru/external/image/" + fmt.Sprint(msgID)
	bts, err := dry.FileGetBytes(imageURI)
	if err != nil {
		color.Red("http %s", err)
		return
	}
	f := tgbotapi.FileBytes{
		Bytes: bts,
		Name:  "file.jpg",
	}
	cnf := tgbotapi.NewPhotoUpload(userID, f)
	cnf.Caption = caption
	cnf.DisableNotification = true
	if len(kbs) > 0 {
		cnf.ReplyMarkup = kbs[0].Markup()
	}
	resp, err := tw.Bot().Send(cnf)
	if err != nil {
		color.Red("send %s", err)
		return
	}
	var maxSize int
	for _, v := range *resp.Photo {
		if maxSize < v.Height {
			maxSize = v.Height
			photoID = v.FileID
		}
	}
	return
}

type Message struct {
	Category int `json:"category"`

	ID int64 `json:"id"`

	GlobalID string `json:"global_id"`
	IsTest   bool   `json:"is_test,omitempty"`

	Contact string `json:"contact"`

	HasPhoto bool   `json:"has_photo"`
	Body     string `json:"body"`
}

type Proposal struct {
	Text   string
	Images []string
}

var proposals = map[int64]Proposal{}

func Add(c *tamework.Context) {
	if len(proposals[c.ChatID].Images) == 0 {
		c.Keyboard = tamework.NewKeyboard([]string{buttons.CancelButton})
	} else {
		c.Keyboard = tamework.NewKeyboard([]string{buttons.CancelButton, buttons.InputText})
	}
	c.Send("Загрузите фотографии (минимум:1, максимум: 10)")

	for {
		proposal, has := proposals[c.ChatID]
		if len(proposal.Images) == 0 {
			has = false
		}
		if has {
			c.Keyboard = tamework.NewKeyboard([]string{buttons.CancelButton, buttons.ClearCurrentProposal, buttons.InputText})
			c.Send(fmt.Sprintf("Вы загрузили %d фотографий", len(proposal.Images)))
		}

		if len(proposal.Images) >= 2 {
			c.Send("Вы загрузили 10 фотографий, теперь введите текст объявления")
		}

		update, done := c.Wait(buttons.CancelButton, time.Second*60)
		color.Cyan("%v %v", update, done)
		if !done {
			color.Green("False on failt %s", update.Method())
			Greeting(c)
			return
		}
		if (update.Message == nil ||
			update.Message.Photo == nil ||
			len(*update.Message.Photo) == 0) &&
			(update.Message == nil ||
				update.Message.Document == nil ||
				update.Message.Document.Thumbnail == nil) {
			if update.Text() != buttons.InputText {
				if update.Text() == buttons.ClearCurrentProposal {
					proposals[c.ChatID] = Proposal{}
					c.Send("Черновик очищен")
					Add(c)
					return
				}
				c.Send("Загрузите фотки")
				continue
			} else {
				c.Send("Введите текст")
				_, done := c.Wait(buttons.CancelButton)
				if !done {
					c.Keyboard = tamework.NewKeyboard(buttons.Menu)
					c.Keyboard.SetRowLen(2)
					c.Send("Возвращайтесь позже, ваш черновик сохранён")
					return
				}

				c.Keyboard = tamework.NewKeyboard([]string{"Предпросмотр", "Опубликовать"})
				c.Keyboard.SetType(tamework.KeyboardInline)
				c.Send("Ваше объявление готово!")
				break
			}

		}

		var maxSizeID string
		var maxSizeValue int

		if update.Message.Document != nil {
			maxSizeID = update.Message.Document.FileID
		} else {
			for _, v := range *update.Message.Photo {
				if v.FileSize > maxSizeValue {
					maxSizeID = v.FileID
					maxSizeValue = v.FileSize
				}
				log.Println(v.FileID, v.FileSize, v.Height, v.Width)
			}
		}

		proposal.Images = append(proposal.Images, maxSizeID)
		proposals[c.ChatID] = proposal
		c.Keyboard = tamework.NewKeyboard([]string{buttons.CancelButton, buttons.InputText})
	}

}

func Terms(c *tamework.Context) {
	c.Keyboard = tamework.NewKeyboard(buttons.Menu)
	c.Keyboard.SetRowLen(2)
	c.Send("Мы не передаём третьим лицам ваши персональные данные, не парьтесь)")
}

func Support(c *tamework.Context) {
	supportText := `По вопросам оплаты и паблика ВК пишите @gnezdovchenko

По вопросам работы бота пишите @zhuha

Если у вас общий вопрос или предложение, пишите тут, мы ответим в течении 8 часов.`

	c.Keyboard = tamework.NewKeyboard(buttons.CancelButton)
	c.Send(supportText)

	u, done := c.Wait(buttons.CancelButton, time.Second*180)
	if !done {
		Greeting(c)
		return
	}
	c.Keyboard = tamework.NewKeyboard(buttons.Menu)
	c.Keyboard.SetRowLen(2)
	c.Send("Мы получили ваш вопрос и уже начали думать)")

	c.Keyboard = tamework.NewKeyboard(nil)
	c.Keyboard.AddCallbackButton("Ответить", "sup"+fmt.Sprint(c.ChatID))
	c.SendTo(setting.App.Telegram.Admin, u.Text())
}

func Greeting(c *tamework.Context) {
	setDefaultKeyboard(c)
	err := c.Markdown(fmt.Sprintf(`Что бы получать новые объявления, выбирайте *%s* и настраивайте нужные параметры.

Что бы добавить объявление, жмите *%s*.`, buttons.SubscriptionsButton, buttons.SubscriptionsOwnerButton))
	if err != nil {
		log.Println(err)
	}
}

func setDefaultKeyboard(c *tamework.Context) {
	c.NewKeyboard(buttons.Menu)
	c.Keyboard.SetRowLen(2)
	c.Keyboard.SetType(tamework.KeyboardReply)
}
