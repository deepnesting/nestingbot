package buttons

var (
	SubscriptionsButton      = "Найти гнёздышко"
	SubscriptionsOwnerButton = "Сдать гнёздышко"
	//AddProposal          = "Сдать гнёздышко"
	AddProposalSearch    = "Подать объявление"
	AboutRent            = "О сдаче квартир"
	AboutRentRoom        = "О сдаче комнат"
	AboutNeightborg      = "О поиске соседа"
	AboutFinderRoom      = "О поиске комнаты"
	AboutFinder          = "О поиске квартиры"
	Back                 = "Назад"
	HelpButton           = "Помощь"
	CancelButton         = "Отменить"
	InputText            = "Ввести текст"
	ClearCurrentProposal = "Очистить черновик"
	FiltersMenu          = []string{
		AboutRent,
		AboutRentRoom,
		AboutNeightborg,
		"",
		//AddProposalSearch,
		"",
		Back,
	}
	FiltersOwnerMenu = []string{
		AboutFinderRoom,
		AboutFinder,
		//AddProposalSearch,
		"",
		Back,
	}
	Menu = []string{
		SubscriptionsButton,
		SubscriptionsOwnerButton,
		//HelpButton,
	}
)
