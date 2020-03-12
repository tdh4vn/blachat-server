package repo

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
	"time"
)

type ChannelRepo interface {
	Create(name string) (*entities.Channel, error)

	FindById(id gocql.UUID) (*entities.Channel, error)

	ChannelsUserJoin(userId gocql.UUID, lastUpdate *time.Time, pageSize int) ([]*entities.Channel, error)

	Delete(id gocql.UUID) (bool, error)

	UpdateLastMessage(id gocql.UUID, mID gocql.UUID) (*entities.Channel, error)
}
