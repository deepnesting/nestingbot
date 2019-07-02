package offers

import (
	"fmt"
	"strings"
)

func FormatMarkdown(o Offer, username string) (text string, err error) {

	if username != "без юзернейма" {
		username = "@" + username
	}
	f := `*#%[1]d* от %[9]s

*%[8]s*, *%[7]s*

Текст: *%[2]s*
Изображения: *%[3]d*
Контакты: *%[6]s*
Тэги: *%[10]s*

Оплачено: *%[4]s*
Опубликовано: *%[5]s*

`
	return strings.TrimSpace(fmt.Sprintf(f,
		o.ID, o.Text, len(o.Images), ifDa(o.Paid), ifDa(o.Published), o.Contacts, city(o.City),
		offerType(o.Type), username,
		strings.Replace(strings.Join(o.Tags, ","), "#", "", -1))), nil
}

func ifDa(da bool) string {
	if da {
		return "да"
	}
	return "нет"

}

func offerType(typ string) string {
	switch typ {
	case SearchNest:
		return "Поиск гнезда"
	case SearchCompanion:
		return "Поиск соседа"
	case SearchNester:
		return "Сдача гнезда"
	case SearchHookUp:
		return "Подселение"
	default:
		return "Ошибка!"
	}
}

func city(id int) string {
	switch id {
	case 1:
		return "Санкт-Петербург"
	case 2:
		return "Москва"
	default:
		return "Не установлен"
	}
}
