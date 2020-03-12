package repo_impl

import (
	"github.com/gocql/gocql"
	"time"
)

const inviteUserToChannelQuery = "INSERT INTO participants(user_id, channel_id, last_receive, last_seen, created_at, updated_at) VALUES(?,?,?,?,?,?)"

const usersInChannelQuery = "SELECT user_id FROM participants WHERE channel_id = ?"

type ParticipantsRepoImpl struct {
	DbSession *gocql.Session
}

func NewParticipantsRepo(db *gocql.Session) *ParticipantsRepoImpl {
	return &ParticipantsRepoImpl{
		DbSession: db,
	}
}

func (p *ParticipantsRepoImpl) InviteToChannel(cID gocql.UUID, userIds []gocql.UUID) error {
	for _, uId := range userIds {
		_ = p.DbSession.Query(inviteUserToChannelQuery, uId, cID, time.Now(), time.Now(), time.Now(), time.Now()).Exec()
	}

	return nil
}

func (p *ParticipantsRepoImpl) UserIdsInChannel(cID gocql.UUID) ([]gocql.UUID, error) {
	var ids []gocql.UUID

	var id string

	iter := p.DbSession.Query(usersInChannelQuery, cID).Iter()

	for iter.Scan(&id) {
		if uuid, err := gocql.ParseUUID(id); err == nil {
			ids = append(ids, uuid)
		} else {
			print(err.Error())
		}
	}

	return ids, nil
}

