package interfaces

import "github.com/DKhorkov/plantsCareTelegramBot/internal/entities"

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	// Users:

	SaveUser(user entities.User) (int, error)
	GetUserByTelegramID(telegramID int) (*entities.User, error)

	// Groups:

	GetUserGroups(userID int) ([]entities.Group, error)
	CountUserGroups(userID int) (int, error)
	AddGroupTitle(telegramID int, title string) (*entities.Group, error)

	// Plants:

	GetUserPlants(userID int) ([]entities.Plant, error)
	CountUserPlants(userID int) (int, error)

	// Temporary:

	GetUserTemporary(telegramID int) (*entities.Temporary, error)
	SetTemporaryStep(telegramID int, step int) error
	SetTemporaryMessage(telegramID int, messageID int) error
}
