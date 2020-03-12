package repo

import (
	"github.com/gocql/gocql"
)

type ParticipantsRepo interface {
	InviteToChannel(cID gocql.UUID, userIds []gocql.UUID) error

	UserIdsInChannel(cID gocql.UUID) ([]gocql.UUID, error)
}