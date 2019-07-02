package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	helpPkg "github.com/deepnesting/nestingbot/domain/help"

	"github.com/Unknwon/com"
	"github.com/asdine/storm"
	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
	offersPkg "github.com/deepnesting/nestingbot/domain/offers"
	"github.com/deepnesting/nestingbot/domain/offers/delivery"
	offersRepoPkg "github.com/deepnesting/nestingbot/domain/offers/repo"
	"github.com/deepnesting/nestingbot/domain/user"
	"github.com/deepnesting/nestingbot/pkg/binlog"
	"github.com/deepnesting/nestingbot/pkg/buttons"
	"github.com/deepnesting/nestingbot/routers/subscriptions"
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhuharev/talert"
	"github.com/zhuharev/tamework"

	whuClient "github.com/zhuharev/whu/domain/client"
)

const version = "0.0.5"

var (
	tw       *tamework.Tamework
	adminIDs []int64
)

// func sendTestMessage() {
// 	bot, _ := tamework.New("882530435:AAF9rkG95tg4f10YNSfOdDFogVYFzgLskwU")
// 	err := bot.Send(132101139, "Добрый день! При создании объявления не сохранились ваши контакты, пришлите их пожалуйста администратору @Andrewsaltanov . Спасибо!")
// 	if err != nil {
// 		panic(err)
// 	}
// }

func getUsername(ctx *tamework.Context) string {
	username := ctx.Update().Username()
	if strings.HasPrefix(username, "_") {
		username = "без юзернейма"
	}
	return username
}

type YaCallback struct {
	Amount           string    `json:"amount"`
	Codepro          string    `json:"codepro"`
	Currency         string    `json:"currency"`
	Datetime         time.Time `json:"datetime"`
	Label            string    `json:"label"`
	NotificationType string    `json:"notification_type"`
	OperationID      string    `json:"operation_id"`
	OperationLabel   string    `json:"operation_label"`
	Sender           string    `json:"sender"`
	Sha1Hash         string    `json:"sha1_hash"`
	TestNotification string    `json:"test_notification"`
}

func makeHandleYACB(bot *tamework.Tamework, offerRepo offersPkg.Repository, userRepo user.Repo) func(data []byte) error {
	return func(data []byte) error {
		log.Info("handle webhook")
		var yaCB YaCallback
		err := json.Unmarshal(data, &yaCB)
		if err != nil {
			return errors.Wrap(err, "unmarshal")
		}
		if !strings.HasPrefix(yaCB.Label, "u:") {
			return nil
		}
		amount, err := strconv.ParseFloat(yaCB.Amount, 64)
		if err != nil {
			return errors.Wrap(err, "parse amount")
		}
		if amount < 37 {
			return errors.Wrap(err, "small imcoming amount")
		}

		arr := strings.Split(yaCB.Label, ":")
		if len(arr) != 2 {
			return nil
		}
		id, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			log.Error("parse label to int", rz.Err(err))
			return nil // may be not bot payment
		}

		offer, err := offerRepo.GetByID(int(id))
		if err != nil {
			return err
		}

		offer.Paid = true

		err = offerRepo.Update(offer)
		if err != nil {
			return err
		}

		user, err := userRepo.Get(offer.UserID)
		if err != nil {
			return err
		}

		// to user
		bot.Send(int64(offer.UserID), "Оплата прошла успешно")

		renderedOffer, _ := offersPkg.FormatMarkdown(*offer, user.Username)
		for _, v := range adminIDs {

			kb := tamework.NewKeyboard(nil).
				AddCallbackButton("Меню", "showmenu:"+arr[1])

			msg := tgbotapi.NewMessage(v, renderedOffer)
			msg.ReplyMarkup = kb.Markup()

			_, err := bot.Bot().Send(msg)
			if err != nil {
				log.Error("send msg to tg", rz.Err(err))
			}

			//bot.Send(v, fmt.Sprintf("Обявление #%d оплачено", offer.ID))
		}

		return nil
	}
}

func main() {
	logger := rz.New(rz.Formatter(rz.FormatterLogfmt()), rz.Level(rz.DebugLevel))
	log.SetLogger(logger)
	//sendTestMessage()
	db, err := storm.Open(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatal("open db")
	}

	userRepo := user.NewRepository(db)
	offerRepo := offersRepoPkg.New(db, userRepo)
	helpRepo := helpPkg.NewRepo(db)

	accessToken := os.Getenv("T_TOKEN")
	adminIDsStr := strings.Split(os.Getenv("ADMIN_IDS"), ",")

	for _, idstr := range adminIDsStr {
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			panic(err)
		}
		adminIDs = append(adminIDs, id)
	}

	err = binlog.NewContext()
	if err != nil {
		log.Fatal("Failed binlog initialization", rz.Err(err))
	}

	tw, err = tamework.New(accessToken)
	if err != nil {
		log.Fatal("Failed telegram bot initialization", rz.Err(err))
	}

	go func() {
		cli := whuClient.New(os.Getenv("WHU_URL"))
		fn := makeHandleYACB(tw, offerRepo, userRepo)
		cli.Run(fn)
	}()
	ttoken, tid, err := talert.ParseDSN(os.Getenv("TALERT_DSN"))
	if err != nil {
		log.Error("parse talert dsn",
			rz.String("dsn", os.Getenv("TALERT_DSN")),
		)
	}
	talert.Init(ttoken, tid)
	talert.Alert("nestingbot started",
		talert.String("version", version))

	tw.NotFound = func(ctx *tamework.Context) {
		ctx.Keyboard.Remove()
		ctx.Send("Вы ввели несуществующую команду!")
		ctx.Keyboard.AddCallbackButton("На главную", "main")
		ctx.Send("Попробуйте начать сначала")
	}

	tw.Bot().RemoveWebhook()

	tw.Use(func(ctx *tamework.Context) {
		if ctx.Text != "" {
			log.Debug("middleware",
				rz.String("text", ctx.Text))
		}
		user, err := userRepo.Get(int(ctx.UserID))
		if err != nil {
			if err == storm.ErrNotFound {
				ctx.Keyboard.Remove()
				ctx.Send("Добро пожаловать!")
				ctx.Keyboard.AddCallbackButton("Санкт-Петербург", "setcity:1")
				ctx.Keyboard.AddCallbackButton("Москва", "setcity:2")
				ctx.Send("Выберите город")
				ctx.Exit()
				userRepo.Create(int(ctx.UserID), "", 0)
				return
			} else {
				log.Error("get user from db", rz.Err(err))
			}
		}

		if user.Username != getUsername(ctx) {
			user.Username = getUsername(ctx)
			err := userRepo.Update(user)
			if err != nil {
				log.Error("update username", rz.Err(err))
			}
		}

		ctx.Data["user"] = user

	})

	tw.Prefix("toggletag:", func(ctx *tamework.Context) {
		log.Debug("toggle tag", rz.String("text", ctx.Text))
		var arr = strings.Split(ctx.Text, ":")
		if len(arr) != 2 {
			return
		}
		id, err := strconv.ParseInt(arr[0], 10, 64)
		if err != nil {
			log.Error("parse id", rz.Err(err))
			return
		}
		offer, err := offerRepo.GetByID(int(id))
		if err != nil {
			log.Error("err get offer from db", rz.Err(err))
			return
		}

		err = offerRepo.ToggleTag(int(id), arr[1])
		if err != nil {
			log.Error("err get offer from db", rz.Err(err))
			return
		}

		//✖️✅
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
		ctx.EditReplyMarkup(ctx.Keyboard)
	})

	tw.Prefix("showmenu:", func(ctx *tamework.Context) {
		//id, _ := strconv.ParseInt(ctx.Text, 10, 64)
		log.Debug("show menu", rz.String("text", ctx.Text))
		ctx.Keyboard.AddCallbackButton("Изображения", "showimages:"+ctx.Text)
		ctx.Keyboard.AddCallbackButton("Опубликовать", "publish:"+ctx.Text)
		ctx.EditReplyMarkup(ctx.Keyboard)
	})

	tw.CallbackQuery("main", func(ctx *tamework.Context) {
		ctx.Answer("Домой")
		setDefaultKeyboard(ctx)
		ctx.Send("Чтобы создать объявление, выберите соответствующие пункты меню")
	})

	tw.Prefix("setcity:", func(ctx *tamework.Context) {
		id, err := strconv.ParseInt(ctx.Text, 10, 64)
		if err != nil {
			log.Error("setcity: parse city id", rz.Err(err), rz.String("text", ctx.Text))
			ctx.Answer("Ошибка установки города")
			return
		}
		if id != 1 && id != 2 {
			log.Error("setcity: bad city id", rz.String("text", ctx.Text), rz.Int64("parsed id", id))
			ctx.Answer("Что-то пошло не так")
			return
		}

		err = userRepo.Create(int(ctx.UserID), getUsername(ctx), int(id))
		if err != nil {
			log.Error("setcity: insert user", rz.Err(err), rz.Int64("tuid", ctx.UserID), rz.Int64("city", id))
			ctx.Answer("Ошибка установки города")
			return
		}
		ctx.Answer("Установлен город: " + ctx.Text)

		setDefaultKeyboard(ctx)
		ctx.Send("Чтобы создать объявление, выберите соответствующие пункты меню")
	})

	tw.Prefix("showhistory:", func(ctx *tamework.Context) {
		id, _ := strconv.ParseInt(ctx.Text, 10, 64)
		list, err := helpRepo.List(int(id))
		if err != nil {
			log.Error("get list from db", rz.Err(err))
		}
		text := ""
		for _, msg := range list {
			if msg.IsIncoming {
				text += "пользователь>" + msg.Text + "\n\n"
			} else {
				text += "адм>" + msg.Text + "\n\n"
			}
		}
		ctx.Keyboard.AddCallbackButton("Ответить", "sup"+strconv.Itoa(int(id)))
		ctx.Send(text)
	})

	tw.Text("/menu", Greeting)
	//alias
	tw.RegistreMethod("/menu", "/start")
	tw.RegistreMethod("/menu", buttons.Back)
	tw.RegistreMethod("/menu", buttons.CancelButton)
	//tw.Text("/start", Greeting)
	tw.Text(buttons.SubscriptionsButton, subscriptions.Subscriptions)
	tw.Text(buttons.SubscriptionsOwnerButton, subscriptions.SubscriptionsOwner)

	tw.Text(buttons.SearchNester, delivery.MakeCreate(offersPkg.SearchNester, adminIDs, offerRepo))
	tw.Text(buttons.SearchCompanion, delivery.MakeCreate(offersPkg.SearchCompanion, adminIDs, offerRepo))
	tw.Text(buttons.SearchHookUp, delivery.MakeCreate(offersPkg.SearchHookUp, adminIDs, offerRepo))
	tw.Text(buttons.SearchNest, delivery.MakeCreate(offersPkg.SearchNest, adminIDs, offerRepo))

	tw.Prefix("Показать объявление ", func(ctx *tamework.Context) {
		id, _ := strconv.Atoi(ctx.Text)
		o, err := offerRepo.GetByID(id)
		if err != nil {
			log.Error("get offers by user", rz.Err(err))
		}
		user, err := userRepo.Get(o.UserID)
		if err != nil {
			log.Error("get user by id", rz.Err(err))
		}
		text, _ := offersPkg.FormatMarkdown(*o, user.Username)
		ctx.Keyboard.AddCallbackButton("Изображения", "showimages:"+strconv.Itoa(int(o.ID)))
		ctx.Keyboard.AddCallbackButton("Файлы", "showfiles:"+strconv.Itoa(int(o.ID)))
		log.Debug("send offer", rz.String("text", text))
		_, err = ctx.Markdown(text)
		if err != nil {
			log.Error("send msg", rz.Err(err))
		}
	})

	tw.Text(buttons.MyOffers, func(ctx *tamework.Context) {
		offrs, err := offerRepo.GetByUserID(int(ctx.UserID))
		if err != nil {
			log.Error("get offers by user", rz.Err(err))
		}
		userFace := ctx.Data["user"]
		user := userFace.(*user.User)
		for _, offer := range offrs {
			text, _ := offersPkg.FormatMarkdown(offer, user.Username)
			ctx.Keyboard.AddCallbackButton("Изображения", "showimages:"+strconv.Itoa(int(offer.ID)))
			log.Debug("send offer", rz.String("text", text))
			_, err = ctx.Markdown(text)
			if err != nil {
				log.Error("send msg", rz.Err(err))
			}
		}
	})

	tw.Prefix("showimages:", func(ctx *tamework.Context) {
		log.Debug("show images", rz.String("text", ctx.Text))
		id, err := strconv.ParseInt(ctx.Text, 10, 64)
		if err != nil {
			log.Error("send msg", rz.Err(err))
			return
		}
		offer, err := offerRepo.GetByID(int(id))
		if err != nil {
			log.Error("send msg", rz.Err(err))
			return
		}
		var images []interface{}
		for _, img := range offer.Images {
			med := tgbotapi.NewInputMediaPhoto(img)
			images = append(images, med)
		}
		msg := tgbotapi.NewMediaGroup(ctx.ChatID, images)
		_, err = ctx.BotAPI().Send(msg)
		if err != nil {
			log.Error("send msg", rz.Err(err))
		}
		ctx.Answer("")
	})

	tw.Prefix("showfiles:", func(ctx *tamework.Context) {
		log.Debug("show images", rz.String("text", ctx.Text))
		id, err := strconv.ParseInt(ctx.Text, 10, 64)
		if err != nil {
			log.Error("send msg", rz.Err(err))
			return
		}
		offer, err := offerRepo.GetByID(int(id))
		if err != nil {
			log.Error("send msg", rz.Err(err))
			return
		}
		for _, img := range offer.Images {
			msg := tgbotapi.NewPhotoShare(ctx.ChatID, img)
			msg.Caption = img
			ctx.BotAPI().Send(msg)
		}
		ctx.Answer("")
	})

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

	tw.Text("/support", MakeSupport(helpRepo))
	tw.CallbackQuery("/support", MakeSupport(helpRepo))
	tw.RegistreMethod("/support", buttons.HelpButton)

	tw.Text(buttons.AddProposalSearch, Add)
	//tw.RegistreMethod(buttons.AddProposal, buttons.AddProposalSearch)

	tw.Text("/terms", Terms)

	// chat
	tw.Prefix("sup", MakeSup(helpRepo))

	//tw.Prefix("setcat_", SetCat)

	// tw.Prefix("upvote_", UpVote)
	// tw.Prefix("downvote_", DownVote)
	tw.Prefix("publish:", MakePublish(offerRepo, userRepo))

	tw.Run()

	// m := macaron.New()
	// m.Post("/cb", func(ctx *macaron.Context) {
	// 	ctx.Req.ParseForm()
	// 	log.Debug("post form", rz.String("form", fmt.Sprint(ctx.Req.PostForm)))

	// 	amount := ctx.QueryFloat64("amount")
	// 	label := ctx.Query("label")
	// 	// todo get offer by label
	// 	// send to moderation
	// 	log.Debug("received values", rz.Float64("amount", amount), rz.String("label", label))
	// })
	// m.Use(macaron.Renderer())

	//m.Post(fmt.Sprintf("/%s/event", setting.App.Secret), binding.Bind(Message{}), eventHandler)

	//m.Run(2018)

}

func MakePublish(offersRepo offersPkg.Repository, userRepo user.Repo) tamework.HandleFunc {
	return func(ctx *tamework.Context) {
		log.Debug("publish offer", rz.String("text", ctx.Text))

		offerID := com.StrTo(ctx.Text).MustInt()

		offer, err := offersRepo.GetByID(offerID)
		if err != nil {
			log.Error("get offer from db", rz.Err(err))
		}
		if offer.Published {
			ctx.Answer("уже опубликовано")
		} else {
			ctx.Answer("опубликовано")
		}

		user, err := userRepo.Get(offer.UserID)
		if err != nil {
			log.Error("get offer user from db", rz.Err(err))
			return
		}

		text := ""
		for _, tag := range offer.Tags {
			text += tag + "\n"
		}
		text += "\n" + offer.Text
		text += "\n\nКонтакты: " + offer.Contacts

		offer.Published = true
		err = offersRepo.Update(offer)
		if err != nil {
			log.Error("set published to db", rz.Err(err))
			return
		}

		reqMsg := tgbotapi.NewMessageToChannel(channelByCity(user.City), text)
		reqMsg.DisableNotification = true
		msg, err := ctx.BotAPI().Send(reqMsg)
		if err != nil {
			log.Error("err send message", rz.Err(err))
			return
		}

		var images []interface{}
		for _, img := range offer.Images {
			med := tgbotapi.NewInputMediaPhoto(img)
			images = append(images, med)
		}
		mdg := tgbotapi.NewMediaGroup(ctx.ChatID, images)
		mdg.ChannelUsername = "@zhutest"
		mdg.ReplyToMessageID = msg.MessageID
		mdg.DisableNotification = true
		_, err = ctx.BotAPI().Send(mdg)
		if err != nil {
			log.Error("send msg", rz.Err(err))
		}
	}
}

// func UpVote(c *tamework.Context) {
// 	vote(c, 1)
// }

// func DownVote(c *tamework.Context) {
// 	vote(c, -1)
// }

// func vote(c *tamework.Context, vote int) {
// 	user, err := models.UserGetOrCreate(c.UserID)
// 	if err != nil {
// 		color.Red("%s", err)
// 		return
// 	}
// 	objID := com.StrTo(c.Text).MustInt64()
// 	err = models.VotesVote(user.ID, objID, vote)
// 	if err != nil {
// 		color.Red("%s", err)
// 		return
// 	}
// 	c.Answer("Голос учтён")
// 	cnt, err := models.VotesCount(objID, 1)
// 	if err != nil {
// 		color.Red("%s", err)
// 		return
// 	}
// 	cntDown, err := models.VotesCount(objID, -1)
// 	if err != nil {
// 		color.Red("%s", err)
// 		return
// 	}
// 	isAdmin := c.ChatID == 102710272
// 	c.EditReplyMarkup(getVoteKeyboard(objID, isAdmin, cnt, cntDown))
// }

// func SetCat(c *tamework.Context) {
// 	arr := strings.Split(c.Text, "_")
// 	if len(arr) != 2 {
// 		return
// 	}
// 	msgID := com.StrTo(arr[0]).MustInt64()
// 	catID := com.StrTo(arr[1]).MustInt()
// 	if cap, ok := catWait[msgID]; ok {
// 		err := broadcast(msgID, catID, cap)
// 		if err != nil {
// 			log.Debug("err broadcast", rz.Err(err))
// 		}
// 		delete(catWait, msgID)
// 	}
// 	c.Answer("Сообщение опубликовано")
// 	c.EditText(c.Update().CallbackQuery.Message.Text)
// }

func MakeSup(helpRepo helpPkg.Repo) tamework.HandleFunc {
	return func(c *tamework.Context) {
		userID := com.StrTo(c.Text).MustInt64()
		c.Keyboard.AddReplyButton(buttons.CancelButton)
		c.Send("Введите ответное сообщение:")
		u, done := c.Wait(buttons.CancelButton, time.Second*180)
		if done {
			c.Keyboard.Reset().AddCallbackButton("Ответить", "/support")
			log.Debug("send help response to user", rz.Int64("user", userID))
			_, err := c.SendTo(userID, "Ответ администратора: \n\n"+u.Text())
			if err != nil {
				log.Error("err send response to message", rz.Err(err))
			}
			helpRepo.Create(int(userID), u.Text(), false)
			setDefaultKeyboard(c)
			c.Send("Сообщение отправлено.")
		}
	}
}

var catWait = map[int64]Message{}

// func eventHandler(c *macaron.Context, event Message) {
// 	if event.IsTest {
// 		c.JSON(200, "ok")
// 		return
// 	}

// 	if event.Category == 0 {
// 		kb := tamework.NewKeyboard(nil).AddCallbackButton("#сдаю_гнездышко", fmt.Sprintf("setcat_%d_1", event.ID)).
// 			AddCallbackButton("").
// 			AddCallbackButton("#сосед_гнездышко", fmt.Sprintf("setcat_%d_3", event.ID)).
// 			AddCallbackButton("").
// 			// AddCallbackButton("О поиске комнат", fmt.Sprintf("setcat_%d_4", event.ID)).
// 			// AddCallbackButton("").
// 			AddCallbackButton("#сниму_гнездышко", fmt.Sprintf("setcat_%d_5", event.ID))
// 		if event.HasPhoto {
// 			for _, admID := range adminIDs {
// 				_, err := uploadRemoteMessage(admID, event.ID, event.Contact, kb)
// 				if err != nil {
// 					log.Debug("err broadcast", rz.Err(err))
// 				}
// 			}
// 		} else {
// 			for _, admID := range adminIDs {
// 				msg := tgbotapi.NewMessage(admID,
// 					fmt.Sprintf("%s\n\n%s", event.Body, event.Contact))
// 				msg.ReplyMarkup = kb.Markup()
// 				_, err := tw.Bot().Send(msg)
// 				if err != nil {
// 					color.Red("%s", err)
// 				}
// 			}
// 		}
// 		catWait[event.ID] = event
// 		return
// 	}
// 	err := broadcast(event.ID, event.Category, event)
// 	if err != nil {
// 		log.Debug("err broadcast", rz.Err(err))
// 	}

// 	c.JSON(200, "ok")
// }

// func getVoteKeyboard(msgID int64, isAdmin bool, votes ...int) *tamework.Keyboard {
// 	kb := tamework.NewKeyboard(nil)
// 	up := "👍"
// 	if len(votes) > 0 && votes[0] != 0 {
// 		up += fmt.Sprintf(" %d", votes[0])
// 	}
// 	down := "👎"
// 	if len(votes) > 1 && votes[1] != 0 {
// 		down += fmt.Sprintf(" %d", votes[1])
// 	}
// 	kb.AddCallbackButton(up, "upvote_"+fmt.Sprint(msgID)).
// 		AddCallbackButton(down, "downvote_"+fmt.Sprint(msgID))
// 	if isAdmin {
// 		kb.AddCallbackButton("в канал "+fmt.Sprint(msgID), "publish_"+fmt.Sprint(msgID))
// 	}
// 	return kb
// }

// func broadcast(msgID int64, catID int, event Message) (err error) {
// 	var (
// 		subs      []int64
// 		channelID = ""
// 	)

// 	switch catID {
// 	//AboutRent
// 	case 1:
// 		channelID = buttons.AboutRent
// 	case 2:
// 		channelID = buttons.AboutRentRoom
// 	case 3:
// 		channelID = buttons.AboutNeightborg
// 	case 4:
// 		channelID = buttons.AboutFinderRoom
// 	case 5:
// 		channelID = buttons.AboutFinder
// 	default:
// 		return
// 	}

// 	subs, err = models.GetSubscribers(channelID)
// 	if err != nil {
// 		return
// 	}

// 	color.Green("%v", subs)
// 	var (
// 		photoID string
// 	)

// 	for _, subID := range subs {
// 		kb := getVoteKeyboard(msgID, false) //

// 		if photoID == "" && event.HasPhoto {
// 			photoID, err = uploadRemoteMessage(subID, msgID, event.Contact, kb)
// 			if err != nil {
// 				color.Red("%s", err)
// 				continue
// 			}
// 			continue
// 		}
// 		if event.HasPhoto {
// 			msg := tgbotapi.NewPhotoShare(subID, photoID)
// 			msg.Caption = event.Contact
// 			msg.ReplyMarkup = kb.Markup()
// 			_, err = tw.Bot().Send(msg)
// 			if err != nil {
// 				color.Red("%s", err)
// 				continue
// 			}
// 		} else {
// 			msg := tgbotapi.NewMessage(subID,
// 				fmt.Sprintf("%s\n\n%s", event.Body, event.Contact))
// 			msg.ReplyMarkup = kb.Markup()
// 			_, err = tw.Bot().Send(msg)
// 			if err != nil {
// 				color.Red("%s", err)
// 				continue
// 			}
// 		}

// 	}
// 	return
// }

// func uploadRemoteMessage(userID int64, msgID int64, caption string, kbs ...*tamework.Keyboard) (photoID string, err error) {
// 	var imageURI = "https://smmpolice.ru/external/image/" + fmt.Sprint(msgID)
// 	bts, err := dry.FileGetBytes(imageURI)
// 	if err != nil {
// 		color.Red("http %s", err)
// 		return
// 	}
// 	f := tgbotapi.FileBytes{
// 		Bytes: bts,
// 		Name:  "file.jpg",
// 	}
// 	cnf := tgbotapi.NewPhotoUpload(userID, f)
// 	cnf.Caption = caption
// 	cnf.DisableNotification = true
// 	if len(kbs) > 0 {
// 		cnf.ReplyMarkup = kbs[0].Markup()
// 	}
// 	resp, err := tw.Bot().Send(cnf)
// 	if err != nil {
// 		color.Red("send %s", err)
// 		return
// 	}
// 	var maxSize int
// 	for _, v := range *resp.Photo {
// 		if maxSize < v.Height {
// 			maxSize = v.Height
// 			photoID = v.FileID
// 		}
// 	}
// 	return
// }

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
				log.Debug("values", rz.String("vals", fmt.Sprintf("%v %v %v %v", v.FileID, v.FileSize, v.Height, v.Width)))
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

func MakeSupport(helpRepo helpPkg.Repo) tamework.HandleFunc {
	return func(c *tamework.Context) {
		supportText := `По вопросам оплаты и паблика ВК пишите @Andrewsaltanov

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

		helpRepo.Create(int(c.UserID), u.Text(), true)

		for _, admID := range adminIDs {
			log.Debug("send help question to admin")
			c.Keyboard = tamework.NewKeyboard(nil)
			username := c.Update().Username()
			if strings.HasPrefix(c.Update().Username(), "_") {
				username = "без юзернейма"
			} else {
				username = "@" + username
			}
			c.Keyboard.AddCallbackButton("Показать переписку", "showhistory:"+fmt.Sprint(c.ChatID))
			c.Keyboard.AddCallbackButton("Ответить", "sup"+strconv.Itoa(int(c.ChatID)))
			msgText := fmt.Sprintf("Новый вопрос (%s):\n\n%s", username, u.Text())
			_, err := c.SendTo(admID, msgText)
			if err != nil {
				log.Error("send msg to t", rz.Err(err), rz.String("text", msgText))
			}
		}
	}
}

func Greeting(c *tamework.Context) {
	setDefaultKeyboard(c)
	_, err := c.Markdown(fmt.Sprintf(`Что бы получать новые объявления, выбирайте *%s* и настраивайте нужные параметры.

Что бы добавить объявление, жмите *%s*.`, buttons.SubscriptionsButton, buttons.SubscriptionsOwnerButton))
	if err != nil {
		log.Error("greeting", rz.Err(err))
	}
}

func setDefaultKeyboard(c *tamework.Context) {
	c.NewKeyboard(buttons.Menu)
	c.Keyboard.SetRowLen(2)
	c.Keyboard.SetType(tamework.KeyboardReply)
}

func channelByCity(cityID int) string {
	switch cityID {
	case 1:
		return "@ughome"
	case 2:
		return "@ugnezdishko"
	default:
		return "@zhutest"
	}
}
