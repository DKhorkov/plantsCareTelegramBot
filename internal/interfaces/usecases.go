package interfaces

import (
	"time"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	// Users:

	SaveUser(user entities.User) (int, error)
	GetUserByTelegramID(telegramID int) (*entities.User, error)

	// Groups:

	GetUserGroups(userID int) ([]entities.Group, error)
	CountUserGroups(userID int) (int, error)
	CreateGroup(group entities.Group) (*entities.Group, error)
	GroupExists(group entities.Group) (bool, error)

	// Plants:

	GetUserPlants(userID int) ([]entities.Plant, error)
	CountUserPlants(userID int) (int, error)

	// Temporary:

	GetUserTemporary(telegramID int) (*entities.Temporary, error)
	SetTemporaryStep(telegramID, step int) error
	SetTemporaryMessage(telegramID, messageID int) error
	AddGroupTitle(telegramID int, title string) (*entities.Group, error)
	AddGroupDescription(telegramID int, description string) (*entities.Group, error)
	AddGroupLastWateringDate(telegramID int, lastWateringDate time.Time) (*entities.Group, error)
	AddGroupWateringInterval(telegramID, wateringInterval int) (*entities.Group, error)
}
