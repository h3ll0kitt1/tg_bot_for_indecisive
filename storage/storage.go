package storage

type Stored interface {
	Save(chatId int, in string) (bool, error)
	Delete(chatId int, in string) (bool, error)
	Exists(chatId int, in string) (bool, error)
	Rand(chatId int) (string, error)
	Print(chatId int) ([]string, error)
	LenNotZero(chatId int) (bool, error)
}
