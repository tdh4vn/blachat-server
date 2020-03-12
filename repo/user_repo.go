package repo

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
)

type UserRepo interface {

	FindByID(id gocql.UUID) (*entities.User, error)

	FindByIDs(ids []gocql.UUID) ([]*entities.User, error)

	Insert(user *entities.User) (*entities.User, error)

	Update(user *entities.User) (*entities.User, error)

}