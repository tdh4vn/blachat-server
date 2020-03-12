package repo

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
)

type MessageRepo interface {
	SaveMessage(message *entities.Message) (*entities.Message, error)
	GetMessages(channelID gocql.UUID, lastItem *gocql.UUID, pageSize int) ([]*entities.Message, error)
	GetNewMessage(channelID gocql.UUID, latestItem *gocql.UUID) ([]*entities.Message, error)
	FindById(messageID gocql.UUID) (*entities.Message, error)
}