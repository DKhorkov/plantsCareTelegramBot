package interfaces

import "github.com/DKhorkov/plantsCareTelegramBot/internal/entities"

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	// Users:

	SaveUser(user entities.User, messageID int) (int, error)

	// Groups:

	GetUserGroups(userID int) ([]entities.Group, error)
	CountUserGroups(userID int) (int, error)

	// Plants:

	GetUserPlants(userID int) ([]entities.Plant, error)
	CountUserPlants(userID int) (int, error)
}
