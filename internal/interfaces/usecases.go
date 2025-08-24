package interfaces

import (
	"time"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	// Users:

	SaveUser(user entities.User) (int, error)
	GetUserByID(id int) (*entities.User, error)
	GetUserByTelegramID(telegramID int) (*entities.User, error)

	// Groups:

	GetUserGroups(userID int) ([]entities.Group, error)
	CountUserGroups(userID int) (int, error)
	CreateGroup(group entities.Group) (*entities.Group, error)
	UpdateGroup(group entities.Group) error
	GroupExists(group entities.Group) (bool, error)
	GetGroup(id int) (*entities.Group, error)
	GetGroupsForNotify(limit, offset int) ([]entities.Group, error)

	// Plants:

	GetUserPlants(userID int) ([]entities.Plant, error)
	CountUserPlants(userID int) (int, error)
	GetGroupPlants(groupID int) ([]entities.Plant, error)
	CountGroupPlants(groupID int) (int, error)
	CreatePlant(plant entities.Plant) (*entities.Plant, error)
	PlantExists(plant entities.Plant) (bool, error)
	GetPlant(id int) (*entities.Plant, error)

	// Temporary:

	GetUserTemporary(telegramID int) (*entities.Temporary, error)
	SetTemporaryStep(telegramID, step int) error
	SetTemporaryMessage(telegramID int, messageID *int) error
	ResetTemporary(telegramID int) error

	AddGroupTitle(telegramID int, title string) (*entities.Group, error)
	AddGroupDescription(telegramID int, description string) (*entities.Group, error)
	AddGroupLastWateringDate(telegramID int, lastWateringDate time.Time) (*entities.Group, error)
	AddGroupWateringInterval(telegramID, wateringInterval int) (*entities.Group, error)

	AddPlantTitle(telegramID int, title string) (*entities.Plant, error)
	AddPlantDescription(telegramID int, description string) (*entities.Plant, error)
	AddPlantGroup(telegramID, groupID int) (*entities.Plant, error)
	AddPlantPhoto(telegramID int, photo []byte) (*entities.Plant, error)
	ManagePlant(telegramID, plantID int) error

	// Notifications:

	SaveNotification(notification entities.Notification) (*entities.Notification, error)
}
