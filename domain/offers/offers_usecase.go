package offers

type UseCase interface {
	CreateOffer(text, images[]string, contacts string) (int,error)
} 