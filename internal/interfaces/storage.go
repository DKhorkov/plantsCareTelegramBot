package interfaces

import "github.com/DKhorkov/plantsCareTelegramBot/internal/entities"

//go:generate mockgen -source=storage.go -destination=../../mocks/storage/storage.go -package=mockstorage
type Storage interface {
	// Users:

	SaveUser(user entities.User) (int, error)
	GetUserByID(id int) (*entities.User, error)
	GetUserByTelegramID(telegramID int) (*entities.User, error)

	// Temporary:

	CreateTemporary(temp entities.Temporary) error
	UpdateTemporary(temp entities.Temporary) error
	GetTemporaryByUserID(userID int) (*entities.Temporary, error)

	// Groups:

	CreateGroup(group entities.Group) (int, error)
	UpdateGroup(group entities.Group) error
	GroupExists(group entities.Group) (bool, error)
	DeleteGroup(id int) error
	GetUserGroups(userID int) ([]entities.Group, error)
	CountUserGroups(userID int) (int, error)
	GetGroup(id int) (*entities.Group, error)
	GetGroupsForNotify(limit, offset int) ([]entities.Group, error)

	// Plants:

	CreatePlant(plant entities.Plant) (int, error)
	UpdatePlant(plant entities.Plant) error
	PlantExists(plant entities.Plant) (bool, error)
	DeletePlant(id int) error
	GetUserPlants(userID int) ([]entities.Plant, error)
	CountUserPlants(userID int) (int, error)
	CountGroupPlants(groupID int) (int, error)
	GetPlant(id int) (*entities.Plant, error)
}
